package config

import (
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	VideoPaths            []string `json:"video_paths"`
	FFMPEGPaths           string   `json:"ffmpeg_paths"`
	FFPROBEPath           string   `json:"ffprobe_path"`
	FileRenameFormatter   string   `json:"file_rename_formatter"`
	FolderRenameFormatter string   `json:"folder_rename_formatter"`
}

func parseJson() *Config {
	fileName := filepath.FromSlash(helpers.GetWorkingDirectory() + "/config.json")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		writeDefault(fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println(err)
		writeDefault(fileName)
	}

	return &config

}

func writeDefault(file string) {
	paths := make([]string, 0)
	paths = append(paths, "G:\\Movies\\tests\\keks\\")
	defaultStruct := Config{VideoPaths: paths,
		FFMPEGPaths:           "",
		FFPROBEPath:           "",
		FileRenameFormatter:   "",
		FolderRenameFormatter: "",
	}

	writeConfig(file, &defaultStruct)
}

func writeConfig(file string, config *Config) {
	Json, _ := json.Marshal(config)
	err := ioutil.WriteFile(file, Json, 0644)

	if err != nil {
		fmt.Println(err)
	}
}

func GetVideoPaths() []string {
	config := parseJson()
	return config.VideoPaths
}

func GetFFMPEGPath() string {
	config := parseJson()
	return config.FFMPEGPaths
}

func GetFFPROBEPath() string {
	config := parseJson()
	return config.FFPROBEPath
}

func WriteFFPROBEPath(path string) {
	config := parseJson()
	config.FFPROBEPath = path
	writeConfig(filepath.FromSlash(helpers.GetWorkingDirectory()+"/config.json"), config)
}

func WriteFFMPEGPath(path string) {
	config := parseJson()
	config.FFMPEGPaths = path
	writeConfig(filepath.FromSlash(helpers.GetWorkingDirectory()+"/config.json"), config)
}

func GetFileRenameFormatter() string {
	config := parseJson()
	return config.FileRenameFormatter
}

func GetFolderRenameFormatter() string {
	config := parseJson()
	return config.FolderRenameFormatter
}

func WriteFileRenameFormatter(formatter string) {
	config := parseJson()
	config.FileRenameFormatter = formatter
	writeConfig(filepath.FromSlash(helpers.GetWorkingDirectory()+"/config.json"), config)
}

func WriteFolderRenameFormatter(formatter string) {
	config := parseJson()
	config.FileRenameFormatter = formatter
	writeConfig(filepath.FromSlash(helpers.GetWorkingDirectory()+"/config.json"), config)
}
