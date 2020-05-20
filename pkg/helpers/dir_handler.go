package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDirs() {
	basePath := GetWorkingDirectory()

	makeDir(basePath + filepath.FromSlash("/assets/database"))
	makeDir(basePath + filepath.FromSlash("/assets/ffmpeg"))
	makeDir(basePath + filepath.FromSlash("/assets/thumbnails"))
	makeDir(basePath + filepath.FromSlash("/logs"))
}

func makeDir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}
