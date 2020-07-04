package ffmpeg

import (
	"bytes"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const component = "FFMPEG"

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
			helpers.LogError(err.Error(), component)
		}
		helpers.LogInfo(fmt.Sprintf("Killed ffmpeg process: %d", cmd.Process.Pid), component)
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
