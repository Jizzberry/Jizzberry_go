package ffmpeg

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/config"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/xi2/xz"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func getUrl() (string, string) {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return "https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-4.1-win64-static.zip", "zip"

		case "386":
			return "https://ffmpeg.zeranoe.com/builds/win32/static/ffmpeg-4.1-win32-static.zip", "zip"
		}
		break

	case "linux":
		switch runtime.GOARCH {
		case "386":
			return "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz", "tar.xz"

		case "arm64":
			return "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz", "tar.xz"
		}
		return "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz", "tar.xz"

	case "darwin":
		return "https://ffmpeg.zeranoe.com/builds/macos64/static/ffmpeg-4.1-macos64-static.zip", "zip"
	}

	return "", ""
}

func untar(path string) error {
	targetDir := filepath.FromSlash(helpers.GetWorkingDirectory() + "/assets/ffmpeg/")

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	r, err := xz.NewReader(f, 0)
	if err != nil {
		return err
	}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(targetDir+hdr.Name, 0777)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			w, err := os.Create(targetDir + hdr.Name)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, tr)
			if err != nil {
				return err
			}
			w.Close()
		}
	}
	f.Close()
	return nil
}

func unzip(path string) error {
	zipReader, _ := zip.OpenReader(path)
	for _, file := range zipReader.Reader.File {

		zippedFile, err := file.Open()
		if err != nil {
			return err
		}

		targetDir := filepath.FromSlash(helpers.GetWorkingDirectory() + "/assets/ffmpeg/")
		extractedFilePath := filepath.Join(
			targetDir,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {

			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return err
			}

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return err
			}
			outputFile.Close()
		}
		zippedFile.Close()
	}
	return nil
}

func DownloadAndExtract() error {

	url, ext := getUrl()
	if url == "" {
		return fmt.Errorf("download ffmpeg manually")
	}

	downloadPath := filepath.FromSlash(helpers.GetWorkingDirectory() + "/assets/ffmpeg/" + "ffmpeg." + ext)

	_ = os.Remove(downloadPath)

	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if ext == "zip" {
		err = unzip(downloadPath)
		if err != nil {
			return err
		}
	}
	if ext == "tar.xz" {
		err = untar(downloadPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func getExecs(path string, file string) string {
	var execPath = ""
	if path == "" {
		path = helpers.GetWorkingDirectory() + "/assets/ffmpeg/"
	}

	if _, err := os.Stat(filepath.FromSlash(path)); err != nil {
		if os.IsNotExist(err) {
			return execPath
		}
	}
	err := filepath.Walk(filepath.FromSlash(path), func(filePath string, f os.FileInfo, err error) error {
		ext := filepath.Ext(filePath)
		if f.IsDir() == false && isValidExt(ext) == true && f.Name() == file+ext {
			execPath = filePath
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	return execPath
}

func isValidExt(ext string) bool {
	switch ext {
	case "",
		".exe":
		return true
	}
	return false
}

func IsExists() error {
	execPathFFMPEG := getExecs(filepath.Dir(config.GetFFMPEGPath()), "ffmpeg")
	execPathProbe := getExecs(filepath.Dir(config.GetFFPROBEPath()), "ffprobe")

	if execPathFFMPEG == "" || execPathProbe == "" {
		err := DownloadAndExtract()

		if err != nil {
			return err
		}

		// Should no longer be empty if download succeeds
		execPathProbe = getExecs("", "ffprobe")
		execPathFFMPEG = getExecs("", "ffmpeg")
	}
	config.WriteFFMPEGPath(execPathFFMPEG)
	config.WriteFFPROBEPath(execPathProbe)
	return nil
}
