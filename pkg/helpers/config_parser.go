package helpers

import (
	"crypto/rand"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	Usernamekey    = "username"
	PasswordKey    = "password"
	SessionsKey    = "sessions"
	LoginURL       = "/auth/login/"
	PrevURLKey     = "prevurl"
	configFileName = "config"

	ThumbnailPath = "./assets/thumbnails"

	component = "Helpers"
)

type Config struct {
	Paths                 []string `json:"paths" mapstructure:"videoPaths"`
	FFMEPG                string   `json:"ffmepg" mapstructure:"ffmpegpath"`
	FFPROBE               string   `json:"ffprobe" mapstructure:"ffprobepath"`
	FileRenameFormatter   string   `json:"file_rename_formatter" mapstructure:"fileRenameFormatter"`
	FolderRenameFormatter string   `json:"folder_rename_formatter" mapstructure:"folderRenameFormatter"`
}

func init() {
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err.Error())
	}

	writeInitial()
}

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

func WriteConfig(config Config) {
	if config.FolderRenameFormatter != "" {
		viper.Set("fileRenameFormatter", config.FolderRenameFormatter)
	}

	if config.FileRenameFormatter != "" {
		viper.Set("fileRenameFormatter", config.FolderRenameFormatter)
	}

	if config.FFPROBE != "" {
		viper.Set("ffprobePath", config.FFPROBE)
	}

	if config.FFMEPG != "" {
		viper.Set("ffmpegPath", config.FFMEPG)
	}

	if config.Paths != nil {
		viper.Set("videoPaths", config.Paths)
	}
	write()
}

func writeInitial() {
	if string(GetSessionsKey()) == "" {
		viper.Set("sessionsKey", GenerateRandomKey(50))
		write()
	}
}

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
	write()
	return nil
}

func RemovePath(path string) error {
	configHolder := GetConfig()
	for i, p := range configHolder.Paths {
		if p == path {
			viper.Set("videoPaths", append(configHolder.Paths[:i], configHolder.Paths[i+1:]...))
			write()
			return nil
		}
	}
	return fmt.Errorf("path not found")
}

func GetSessionsKey() []byte {
	return []byte(viper.GetString("sessionsKey"))
}

func write() {
	if err := viper.WriteConfigAs(configFileName + ".yaml"); err != nil {
		fmt.Println(err)
	}
}

func GenerateRandomKey(l int) string {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		LogError(err.Error(), component)
	}
	return string(b)
}
