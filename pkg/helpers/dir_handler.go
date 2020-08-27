package helpers

import (
	"os"
	"path/filepath"
)

func CreateDirs() error {
	err := makeDir(filepath.FromSlash(DatabasePath))
	if err != nil {
		return err
	}
	err = makeDir(filepath.FromSlash(JsonPath))
	if err != nil {
		return err
	}
	err = makeDir(filepath.FromSlash(FFMPEGPath))
	if err != nil {
		return err
	}
	err = makeDir(filepath.FromSlash(ThumbnailPath))
	if err != nil {
		return err
	}
	err = makeDir(filepath.FromSlash(LogsPath))
	return nil
}

func makeDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm) // Create dir even if parent dir doesn't exist
	if err != nil {
		return err
	}
	return nil
}
