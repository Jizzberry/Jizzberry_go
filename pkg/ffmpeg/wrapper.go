package ffmpeg

import (
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GenerateThumbnail(path string, interval int64, outFile string) {
	var args []string

	// Input fast seeking: https://trac.ffmpeg.org/wiki/Seeking
	args = append(args, "-ss", strconv.FormatInt(interval, 10), "-i", path, "-y", "-vframes:v", "1", "-vf",
		"scale=373:210", filepath.Join(helpers.ThumbnailPath, outFile))

	NewEncoder().Run(args, ThumbnailGen, 10, false)
}

func getFFprobeJson(filepath string) map[string]interface{} {
	var args []string
	args = append(args, "-print_format", "json", "-show_format", "-show_streams", "-v", "quiet", filepath)

	stdout, _, _ := NewEncoder().Run(args, FFProbe, TimeoutBlock, true)
	return parseJson(stdout.Bytes())
}

func Transcode(filepath string, startTime string, encoder *Encoder) (io.ReadCloser, string) {
	var args []string

	if startTime != "" {
		args = append(args, "-ss", startTime)
	}

	args = append(args,
		"-i", filepath,
		"-c:v", "libvpx-vp9",
		"-cpu-used", "6",
		"-deadline", "realtime",
		"-preset", "veryfast",
		"-row-mt", "1",
		"-crf", "30",
		"-b:v", "0",
		"-f", "webm",
		"pipe:",
	)

	stdout, uid := encoder.Pipe(args, TranscodeStream, TimeoutForget)

	return stdout, uid
}

func getLength(data map[string]interface{}) float64 {
	if data != nil {
		duration := helpers.SafeConvertFloat(helpers.SafeSelectFromMap(helpers.SafeMapCast(helpers.SafeSelectFromMap(data, "format")), "duration"))
		return duration
	}
	return -1
}

func ProbeVideo(path string) (length float64, format string, videoCodec string, audioCodec string) {
	data := getFFprobeJson(path)

	format = getVideoFormat(data)
	length = getLength(data)
	videoCodec = getVideoCodec(data)
	audioCodec = getAudioCodec(data)
	return
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

func ConvertFaststart(filepath string) *os.File {
	tmpfile, err := ioutil.TempFile(os.TempDir(), "tmp.*.mp4")
	if err != nil {
		helpers.LogError(err.Error())
		return nil
	}
	var args []string

	// Brings moov atom to start of file
	args = append(args,
		"-i", filepath,
		"-c:a", "copy",
		"-c:v", "copy",
		"-movflags", "faststart",
		"-y",
		tmpfile.Name(),
	)

	NewEncoder().Pipe(args, TranscodeStream, TimeoutBlock)

	return tmpfile

}

func avc1ToRfc6381(tmp *os.File) string {
	buffer := make([]byte, 8128)
	_, err := tmp.Read(buffer)
	if err != nil {
		helpers.LogError(err.Error())
		return ""
	}

	var codecFlags []byte
	// According to http://xhelmboyx.tripod.com/formats/mp4-layout.txt
	r := regexp.MustCompile(string([]byte{97, 118, 99, 67})) // Byte code for "avcC"
	match := r.FindStringIndex(string(buffer))               // Find index for avcC box
	if len(match) > 1 {
		codecFlags = buffer[match[1]+1 : match[1]+4] // Profile, compatibility and level flags
	}

	if len(codecFlags) > 0 {
		return fmt.Sprintf("avc1.%x", codecFlags)
	}
	return ""
}

func GetCodecs(filePath string) string {
	codecs := make([]string, 0)
	tmp := ConvertFaststart(filePath)
	if val := avc1ToRfc6381(tmp); val != "" {
		codecs = append(codecs, val)
	}

	defer func() {
		err := os.Remove(tmp.Name())
		if err != nil {
			helpers.LogError(err.Error())
		}
	}()

	if len(codecs) > 0 {
		return fmt.Sprintf("codecs=\"%s\"", strings.Join(codecs, ","))
	}
	return ""
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
