package ffmpeg

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func GenerateThumbnail(path string, interval int64, outFile string) {
	outPath := filepath.Join(helpers.ThumbnailPath, outFile)
	cmd := exec.Command(helpers.GetConfig().FFMEPG, "-i", path, "-ss", strconv.FormatInt(interval, 10), "-y", "-vframes", "1", "-vf",
		"scale=373:210", outPath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		_ = os.Remove(outPath)
		return
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err = <-done:
		return
	case <-time.After(120 * time.Second):
		err := cmd.Process.Kill()
		if err != nil {
			helpers.LogError(err.Error())
		}
		helpers.LogInfo(fmt.Sprintf("Killed ffmpeg process: %d", cmd.Process.Pid))
		return
	}
}

func getLength(data map[string]interface{}) float64 {
	if data != nil {
		duration := helpers.SafeConvertFloat(helpers.SafeSelectFromMap(helpers.SafeMapCast(helpers.SafeSelectFromMap(data, "format")), "duration"))
		return duration
	}
	return -1
}

func getFFprobeJson(filepath string) map[string]interface{} {
	cmd := exec.Command(helpers.GetConfig().FFPROBE, "-print_format", "json", "-show_format", "-show_streams", "-v", "quiet", filepath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		helpers.LogError("Couldn't execute ffprobe")
		return nil
	}
	return parseJson(out.Bytes())
}

func getAudioCodec(data map[string]interface{}) string {
	streams := helpers.SafeCastSlice(helpers.SafeSelectFromMap(data, "streams"))
	if len(streams) >= 2 {
		return helpers.SafeCastString(helpers.SafeSelectFromMap(helpers.SafeMapCast(streams[1]), "codec_name"))
	}
	return ""
}

func getVideoCodec(data map[string]interface{}) string {
	streams := helpers.SafeCastSlice(helpers.SafeSelectFromMap(data, "streams"))
	if len(streams) >= 1 {
		return helpers.SafeCastString(helpers.SafeSelectFromMap(helpers.SafeMapCast(streams[0]), "codec_name"))
	}
	return ""
}

func getVideoFormat(data map[string]interface{}) string {
	format := helpers.SafeMapCast(helpers.SafeSelectFromMap(data, "format"))
	if val, ok := format["format_name"]; ok {
		return helpers.SafeCastString(val)
	}
	return ""
}

func ProbeVideo(path string) (length float64, format string, videoCodec string, audioCodec string) {
	data := getFFprobeJson(path)

	fmt.Println(data)

	file, err := os.Open(path)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}
	buffer := make([]byte, 3)
	_, err = file.ReadAt(buffer, 539)
	if err != nil {
		helpers.LogError(err.Error())
	}
	fmt.Println(hex.EncodeToString(buffer))
	//fmt.Println(fmt.Sprintf("%x ",buffer))

	format = getVideoFormat(data)
	length = getLength(data)
	videoCodec = getVideoCodec(data)
	audioCodec = getAudioCodec(data)
	return
}

func parseJson(data []byte) map[string]interface{} {
	jsonData := make(map[string]interface{})
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		helpers.LogError(err.Error())
		return nil
	}
	return jsonData
}

func Transcode(filepath string, startTime string) io.ReadCloser {
	var args []string

	if startTime != "" {
		fmt.Println(startTime)
		args = append(args, "-ss", startTime)
	}

	args = append(args,
		"-i", filepath,
		"-c:v", "libvpx-vp9",
		"-vf", "scale=360:-2",
		"-deadline", "realtime",
		"-cpu-used", "5",
		"-row-mt", "1",
		"-crf", "30",
		"-b:v", "0",
		"-f", "webm",
		"pipe:",
	)

	fmt.Println(args)

	cmd := exec.Command(helpers.GetConfig().FFMEPG, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		helpers.LogError(err.Error())
	}

	err = cmd.Start()
	if err != nil {
		helpers.LogError(err.Error())
	}

	return stdout
}

//func ConvertFaststart(filepath string) *os.File {
//	tmpfile, err := ioutil.TempFile(os.TempDir(), "tmp.*.mp4")
//	if err != nil {
//		helpers.LogError(err.Error())
//		return nil
//	}
//	var args []string
//
//	// Brings moov atom to start of file
//	args = append(args,
//		"-i", filepath,
//		"-c:a", "copy",
//		"-c:v", "copy",
//		"-movflags", "faststart",
//		"-y",
//		tmpfile.Name(),
//	)
//	cmd := exec.Command(helpers.GetConfig().FFMEPG, args...)
//	var stderr bytes.Buffer
//	cmd.Stderr = &stderr
//
//
//	err = cmd.Run()
//	if err != nil {
//		helpers.LogError(err.Error())
//	}
//
//	return tmpfile
//
//}
