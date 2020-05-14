package helpers

import (
	"os"
	"path/filepath"
)

func CreateDirs() {
	basePath := GetWorkingDirectory()

	os.MkdirAll(basePath+filepath.FromSlash("/assets/database"), os.ModePerm)
	os.MkdirAll(basePath+filepath.FromSlash("/assets/ffmpeg"), os.ModePerm)
	os.MkdirAll(basePath+filepath.FromSlash("/assets/thumbnails"), os.ModePerm)
}
