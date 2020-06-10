package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDirs() {
	basePath := GetWorkingDirectory()

	makeDir(filepath.Join(basePath, filepath.FromSlash("/assets/database")))
	makeDir(filepath.Join(basePath, filepath.FromSlash("/assets/json")))
	makeDir(filepath.Join(basePath, filepath.FromSlash("/assets/ffmpeg")))
	makeDir(filepath.Join(basePath, filepath.FromSlash(ThumbnailPath)))
	makeDir(filepath.Join(basePath, filepath.FromSlash("/logs")))
}

func makeDir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}
