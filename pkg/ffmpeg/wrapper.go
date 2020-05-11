package ffmpeg

import (
	"bytes"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/config"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GenerateThumbnail(generatedId int64, path string, interval int64) {
	outFile := filepath.FromSlash(helpers.GetWorkingDirectory() + "/assets/thumbnails/" + strconv.FormatInt(generatedId, 10) + ".png")
	cmd := exec.Command(config.GetFFMPEGPath(), "-i", path, "-ss", strconv.FormatInt(interval, 10), "-y", "-vframes", "1", "-vf",
		"scale=540:-1", outFile)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Remove(outFile)
	}
}

func GetLength(path string) string {
	cmd := exec.Command(config.GetFFPROBEPath(), "-v", "error", "-show_entries", "format=duration", "-sexagesimal", "-of", "default=noprint_wrappers=1:nokey=1", path)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

	return strings.Split(out.String(), ".")[0]
}
