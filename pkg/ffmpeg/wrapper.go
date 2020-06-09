package ffmpeg

import (
	"bytes"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const component = "FFMPEG"

func GenerateThumbnail(generatedId int64, path string, interval int64) {
	outFile := filepath.FromSlash(helpers.ThumbnailPath + "/" + strconv.FormatInt(generatedId, 10) + ".png")
	cmd := exec.Command(helpers.GetConfig().FFMEPG, "-i", path, "-ss", strconv.FormatInt(interval, 10), "-y", "-vframes", "1", "-vf",
		"scale=540:-1", outFile)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		_ = os.Remove(outFile)
	}
}

func GetLength(path string) string {
	cmd := exec.Command(helpers.GetConfig().FFPROBE, "-v", "error", "-show_entries", "format=duration", "-sexagesimal", "-of", "default=noprint_wrappers=1:nokey=1", path)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		helpers.LogError(fmt.Sprint(err)+": "+strings.TrimSpace(stderr.String()), component)
	}

	return strings.Split(out.String(), ".")[0]
}
