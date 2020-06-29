package ffmpeg

import (
	"bytes"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const component = "FFMPEG"

func GenerateThumbnail(generatedId int64, path string, interval int64) {
	outFile := helpers.GetThumbnailPath(generatedId, false)
	cmd := exec.Command(helpers.GetConfig().FFMEPG, "-i", path, "-ss", strconv.FormatInt(interval, 10), "-y", "-vframes", "1", "-vf",
		"scale=373:210", outFile)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		_ = os.Remove(outFile)
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
			helpers.LogError(err.Error(), component)
		}
		return
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
