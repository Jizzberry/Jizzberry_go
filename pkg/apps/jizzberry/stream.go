package jizzberry

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/files"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const BUFSIZE = 1024 * 8

func streamHandler(w http.ResponseWriter, r *http.Request) {
	model := files.Initialize()
	defer model.Close()

	vars := mux.Vars(r)
	queryParams := r.URL.Query()

	sceneId, _ := strconv.ParseInt(vars["scene_id"], 10, 64)
	if sceneId == 0 {
		return
	}

	file := model.Get(files.Files{GeneratedID: sceneId})
	if len(file) < 1 {
		return
	}

	path := file[0].FilePath

	if path == "" {
		return
	}

	var playable = false
	if len(queryParams["playable"]) > 0 {
		playable = strings.ToLower(queryParams["playable"][0]) == "true"
	}

	fmt.Println(playable)

	startTime := ""
	if len(queryParams["start"]) > 0 {
		startTime = queryParams["start"][0]
	}

	if playable && IsSupportedCodec(file[0].Video0Codec) {
		pseudoStream(w, r, path, file[0])
	} else {
		transcodeAndStream(w, r, path, startTime)
	}
}

func transcodeAndStream(w http.ResponseWriter, r *http.Request, path string, startTime string) {
	stream := ffmpeg.Transcode(path, startTime)
	_, err := io.Copy(w, stream)
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func pseudoStream(w http.ResponseWriter, r *http.Request, path string, model files.Files) {
	file, err := os.Open(path)
	if err != nil {
		helpers.LogError(err.Error())
		w.WriteHeader(500)
		return
	}

	defer file.Close()

	fi, err := file.Stat()

	if err != nil {
		w.WriteHeader(500)
		return
	}

	fileSize := int(fi.Size())

	if len(r.Header.Get("Range")) == 0 {

		contentLength := strconv.Itoa(fileSize)
		contentEnd := strconv.Itoa(fileSize - 1)

		fmt.Println(MimeTypeFromFormat(model))
		w.Header().Set("Content-Type", MimeTypeFromFormat(model))
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
			_, _ = w.Write(data)
			w.(http.Flusher).Flush()
		}
	}
}
