package api

import (
	"path/filepath"
	"strings"
)

func GetDirectory(path string) (b []browse) {
	if path == "" {
		path = "/"
	}

	allFiles, err := getAllFolders(filepath.FromSlash(path))
	if err != nil {
		return
	}

	for _, f := range allFiles {
		split := strings.Split(f, "/")
		b = append(b, browse{
			Name: split[len(split)-1],
			Path: f,
		})
	}

	b = append(b, browse{
		Name: "..",
		Path: filepath.Join(path, ".."),
	})
	return
}
