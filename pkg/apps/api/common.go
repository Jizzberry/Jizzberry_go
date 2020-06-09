package api

import (
	"os"
	"path/filepath"
)

type browse struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func getAllFolders(path string) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.Walk(filepath.FromSlash(path), func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() == true && filePath != filepath.FromSlash(path) {
			fileList = append(fileList, filePath)
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return fileList, err
	}
	return fileList, nil
}
