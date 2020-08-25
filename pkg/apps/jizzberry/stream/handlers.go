package stream

import (
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/files"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const baseURL = "/stream"

type Stream struct {
	SceneID   int64
	Playable  bool
	Encoder   *ffmpeg.Encoder
	startTime string
}

type StreamResp struct {
	URL      string
	MimeType string
}
type Streamer struct{}

var mapMutex sync.Mutex
var streamHolder = make(map[string]Stream)

const BUFSIZE = 1024 * 8

func (a Streamer) Register(r *mux.Router) {
	streamRouter := r.PathPrefix(baseURL).Subrouter()
	streamRouter.StrictSlash(true)

	streamRouter.HandleFunc("/getStream", newURL)
	streamRouter.HandleFunc("/stream/{uid}", streamHandler)
}

func newURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	queryParams := r.URL.Query()
	if val, ok := queryParams["scene_id"]; ok {
		startTime := "0"
		if len(queryParams["start"]) > 0 {
			startTime = queryParams["start"][0]
		}

		var playable = true
		if len(queryParams["playable"]) > 0 {
			playable = strings.ToLower(queryParams["playable"][0]) == "true"
		}

		sceneId, _ := strconv.ParseInt(val[0], 10, 64)
		resp := URLGenerator(sceneId, playable, startTime)

		jsonEncoder := json.NewEncoder(w)
		err := jsonEncoder.Encode(&resp)
		if err != nil {
			helpers.LogError(err.Error())
		}
	}
}

func URLGenerator(sceneId int64, playable bool, startTime string) StreamResp {
	helpers.LogInfo(sceneId, playable, startTime)

	streamEncoder := ffmpeg.NewEncoder()

	mapMutex.Lock()
	streamHolder[streamEncoder.UID] = Stream{
		Playable:  playable,
		Encoder:   streamEncoder,
		SceneID:   sceneId,
		startTime: startTime,
	}
	mapMutex.Unlock()

	return StreamResp{
		URL:      fmt.Sprintf("/Jizzberry/stream/stream/%s", streamEncoder.UID),
		MimeType: MimeTypeFromFormat(sceneId),
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["uid"]
	if dets, ok := streamHolder[uid]; ok {
		model := files.Initialize()
		defer model.Close()

		file := model.Get(files.Files{GeneratedID: dets.SceneID})

		if len(file) > 0 {
			if dets.Playable {
				pseudoStream(w, r, file[0])
			} else {
				transcodeAndStream(w, dets.Encoder, file[0].FilePath, dets.startTime)
				notify := r.Context().Done()
				go func() {
					<-notify
					dets.Encoder.KillPrev()
				}()
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func transcodeAndStream(w http.ResponseWriter, encoder *ffmpeg.Encoder, path string, startTime string) {
	stream := ffmpeg.Transcode(path, startTime, encoder)
	_, err := io.Copy(w, stream)
	if err != nil {
		helpers.LogWarning(err.Error())
	}
}

func pseudoStream(w http.ResponseWriter, r *http.Request, model files.Files) {
	file, err := os.Open(model.FilePath)
	if err != nil {
		helpers.LogError(err.Error())
		w.WriteHeader(500)
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			helpers.LogError(err.Error())
		}
	}()

	fi, err := file.Stat()

	if err != nil {
		w.WriteHeader(500)
		return
	}

	fileSize := int(fi.Size())

	if len(r.Header.Get("Range")) == 0 {

		contentLength := strconv.Itoa(fileSize)
		contentEnd := strconv.Itoa(fileSize - 1)

		w.Header().Set("Content-Type", MimeTypeFromFormat(model.GeneratedID))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes 0-"+contentEnd+"/"+contentLength)
		w.WriteHeader(200)

		buffer := make([]byte, BUFSIZE)

		for {
			n, err := file.Read(buffer)

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			data := buffer[:n]
			_, _ = w.Write(data)
			w.(http.Flusher).Flush()
		}

	} else {

		rangeParam := strings.Split(r.Header.Get("Range"), "=")[1]
		splitParams := strings.Split(rangeParam, "-")

		// response values
		contentStartValue := 0
		contentStart := strconv.Itoa(contentStartValue)
		contentEndValue := fileSize - 1
		contentEnd := strconv.Itoa(contentEndValue)
		contentSize := strconv.Itoa(fileSize)

		if len(splitParams) > 0 {
			contentStartValue, err = strconv.Atoi(splitParams[0])

			if err != nil {
				contentStartValue = 0
			}

			contentStart = strconv.Itoa(contentStartValue)
		}

		if len(splitParams) > 1 {
			contentEndValue, err = strconv.Atoi(splitParams[1])

			if err != nil {
				contentEndValue = fileSize - 1
			}

			contentEnd = strconv.Itoa(contentEndValue)
		}

		contentLength := strconv.Itoa(contentEndValue - contentStartValue + 1)

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes "+contentStart+"-"+contentEnd+"/"+contentSize)
		w.WriteHeader(206)

		buffer := make([]byte, BUFSIZE)

		_, err := file.Seek(int64(contentStartValue), 0)
		if err != nil {
			helpers.LogError(err.Error())
		}

		writeBytes := 0

		for {
			n, err := file.Read(buffer)

			writeBytes += n

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			if writeBytes >= contentEndValue {
				data := buffer[:BUFSIZE-writeBytes+contentEndValue+1]
				_, _ = w.Write(data)
				w.(http.Flusher).Flush()
				break
			}

			data := buffer[:n]
			_, err = w.Write(data)
			if err != nil {
				helpers.LogError(err.Error())
			}
			w.(http.Flusher).Flush()
		}
	}
}
