package helpers

import (
	"os"
	"path/filepath"
)

var basePath = ""

func GetWorkingDirectory() string {
	return basePath
}

func GetExecutableDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		LogError(err.Error())
	}
	return dir
}

func SetWorkingDirectory(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		LogError(err.Error())
		abs = "."
	}
	basePath = abs
}
