package helpers

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

func IsFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetThumbnailPath() (path string) {
	id := uuid.New()
	path = fmt.Sprintf("%s.png", id.String())
	return
}
