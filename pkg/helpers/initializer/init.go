package initializer

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/apps"
	"github.com/Jizzberry/Jizzberry_go/pkg/database"
	"github.com/Jizzberry/Jizzberry_go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/auth"
	"github.com/Jizzberry/Jizzberry_go/pkg/scrapers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"strings"
	"syscall"
)

// Initialize the whole app
func Init() error {
	err := initHelpers()
	if err != nil {
		return err
	}

	err = database.RunMigrations()
	if err != nil {
		return err
	}

	err = ffmpeg.IsExists()
	if err != nil {
		return err
	}

	err = IsFirstTime()
	if err != nil {
		return err
	}

	scrapers.RegisterScrapers()

	//fmt.Println(ffmpeg.ProbeVideo(files.Initialize().Get(files.Files{})[19].FilePath))
	//fmt.Println(jizzberry.Avc1ToRfc6381(files.Initialize().Get(files.Files{})[20].FilePath))
	//ffmpeg.GetMoovAtom(files.Initialize().Get(files.Files{})[19].FilePath)

	err = initWebApp()
	if err != nil {
		return err
	}
	return nil
}

// Initialize helpers package (Logger, Config, Dirs...)
func initHelpers() error {
	err := helpers.ConfigInit()
	if err != nil {
		return err
	}
	err = helpers.CreateDirs()
	if err != nil {
		return err
	}
	helpers.LoggerInit()
	helpers.RndInit()
	return nil
}

// Initialize Web server (default :6969)
func initWebApp() error {
	fmt.Println(helpers.Art)

	addr := flag.String("addr", ":6969", "Address of server [default :6969]")
	flag.Parse()

	router := mux.NewRouter()

	apps.RegisterFileServer(router)
	apps.RegisterApps(router)

	helpers.LogInfo("Server starting at " + *addr)

	err := http.ListenAndServe(*addr, router)
	if err != nil {
		return err
	}

	return nil
}

// Runs user creation if Auth model is empty
func IsFirstTime() error {
	model := auth.Initialize()
	defer model.Close()

	if len(model.Get(auth.Auth{})) == 0 {
		err := CreateFirstUser(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates user
func CreateFirstUser(model *auth.Model) error {
	username, password, err := inputCreds()
	if err != nil {
		return err
	}

	model.Create(auth.Auth{
		Username: username,
		Password: password,
		IsAdmin:  true,
	})
	return nil
}

// CLI interface to accept credentials
func inputCreds() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter Username: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}

		name = strings.TrimSpace(name)

		fmt.Println("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", err
		}
		password := string(bytePassword)

		fmt.Println("Confirm Password: ")
		bytePasswordC, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", err
		}
		passwordC := string(bytePasswordC)

		if password == passwordC {
			password = strings.Trim(password, "\n")
			return name, password, nil
		} else {
			fmt.Println("Passwords don't match")
		}
	}
}
