package helpers

import (
	"crypto/rand"
	"fmt"
	"github.com/spf13/viper"
	"os"
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

var configFilePath = GetWorkingDirectory() + "/" + configFileName

func init() {
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	writeInitial()

	err := viper.ReadInConfig() // Find and read the helpers file
	if err != nil {             // Handle errors reading the helpers file
		fmt.Print(err.Error())
	}
}

func writeInitial() {
	if string(GetSessionsKey()) == "" {
		viper.Set("sessionsKey", GenerateRandomKey(50))
	}

	write()
}

func GetSessionsKey() []byte {
	return []byte(viper.GetString("sessionsKey"))
}

func GetVideoPaths() []string {
	return viper.GetStringSlice("videoPaths")
}

func GetFFMPEGPath() string {
	return viper.GetString("ffmpegPath")
}

func WriteFFMPEGPath(path string) {
	viper.Set("ffmpegPath", path)
	write()
}

func GetFFPROBEPath() string {
	return viper.GetString("ffprobePath")
}

func WriteFFPROBEPath(path string) {
	viper.Set("ffprobePath", path)
	write()
}

func GetFileRenameFormatter() string {
	return viper.GetString("fileRenameFormatter")
}

func WriteFileRenameFormatter(path string) {
	viper.Set("fileRenameFormatter", path)
	write()
}

func GetFolderRenameFormatter() string {
	return viper.GetString("folderRenameFormatter")
}

func WriteFolderRenameFormatter(path string) {
	viper.Set("folderRenameFormatter", path)
	write()
}

func write() {
	if err := viper.SafeWriteConfigAs(configFilePath); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfigAs(configFilePath)
		}
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