package helpers

import (
	"os"
	"path/filepath"
)

func CreateDirs() error {
	basePath := GetWorkingDirectory()

	err := makeDir(filepath.Join(basePath, filepath.FromSlash("/assets/database")))
	if err != nil {
		return err
	}
	err = makeDir(filepath.Join(basePath, filepath.FromSlash("/assets/json")))
	if err != nil {
		return err
	}
	err = makeDir(filepath.Join(basePath, filepath.FromSlash("/assets/ffmpeg")))
	if err != nil {
		return err
	}
	err = makeDir(filepath.Join(basePath, filepath.FromSlash(ThumbnailPath)))
	if err != nil {
		return err
	}
	err = makeDir(filepath.Join(basePath, filepath.FromSlash("/logs")))
	return nil
}

func makeDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
