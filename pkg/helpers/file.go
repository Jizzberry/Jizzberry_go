package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetThumbnailPath(id int64, actor bool) (path string) {
	if actor {
		path = filepath.Join(ThumbnailPath, fmt.Sprintf("p%d.png", id))
	} else {
		path = filepath.Join(ThumbnailPath, fmt.Sprintf("%d.png", id))
	}
	return
}
