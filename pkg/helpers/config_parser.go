package helpers

import (
	"crypto/rand"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	configFileName = "config"
	configFormat   = "yaml"
)

type Config struct {
	Paths                 []string `json:"paths" mapstructure:"videoPaths"`
	FFMEPG                string   `json:"ffmepg" mapstructure:"ffmpegpath"`
	FFPROBE               string   `json:"ffprobe" mapstructure:"ffprobepath"`
	FileRenameFormatter   string   `json:"file_rename_formatter" mapstructure:"fileRenameFormatter"`
	FolderRenameFormatter string   `json:"folder_rename_formatter" mapstructure:"folderRenameFormatter"`
}

// Initialize Config at provided path
func ConfigInit() error {
	initPaths()

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFormat)
	viper.AddConfigPath(configPath)

	_ = viper.ReadInConfig()

	err := writeInitial()
	if err != nil {
		return err
	}
	return nil
}

// Parses Config doc
func parseConfig() Config {
	var tmp Config
	err := viper.Unmarshal(&tmp)
	if err != nil {
		fmt.Println(err)
	}
	return tmp
}

func GetConfig() Config {
	return parseConfig()
}

// Writes Config to Config file
// ffmpegPath, ffprobePath shouldn't be empty
// videoPaths can be empty but not nil
func WriteConfig(config Config) error {
	viper.Set("folderRenameFormatter", config.FolderRenameFormatter)
	viper.Set("fileRenameFormatter", config.FileRenameFormatter)

	if config.FFPROBE != "" {
		viper.Set("ffprobePath", config.FFPROBE)
	}

	if config.FFMEPG != "" {
		viper.Set("ffmpegPath", config.FFMEPG)
	}

	if config.Paths != nil {
		viper.Set("videoPaths", config.Paths)
	}
	err := write()
	return err
}

// Writes default values if empty
func writeInitial() error {
	if string(GetSessionsKey()) == "" {
		viper.Set("sessionsKey", GenerateRandomKey(50))
		err := write()
		return err
	}
	return nil
}

// Add string to videoPaths if exists as a directory or valid path
func AddPath(path string) error {
	configHolder := GetConfig()
	for _, p := range configHolder.Paths {
		if filepath.FromSlash(p) == filepath.FromSlash(path) {
			return fmt.Errorf("path already exists")
		}
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("invalid path")
		}
	}
	viper.Set("videoPaths", append(configHolder.Paths, path))
	err := write()
	if err != nil {
		return err
	}
	return nil
}

// Removes string from videoPaths if exists
func RemovePath(path string) error {
	configHolder := GetConfig()
	for i, p := range configHolder.Paths {
		if p == path {
			viper.Set("videoPaths", append(configHolder.Paths[:i], configHolder.Paths[i+1:]...))
			err := write()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetSessionsKey() []byte {
	return []byte(viper.GetString("sessionsKey"))
}

func write() error {
	if err := viper.WriteConfigAs(filepath.Join(configPath, configFileName+"."+configFormat)); err != nil {
		return err
	}
	return nil
}

func GenerateRandomKey(l int) string {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		LogError(err.Error())
	}
	return string(b)
}
