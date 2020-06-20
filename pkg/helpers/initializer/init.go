package initializer

import (
	"bufio"
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

var addr = ":8000"

const art = ` 
$$$$$\ $$\                     $$\                                               
   \__$$ |\__|                    $$ |                                              
      $$ |$$\ $$$$$$$$\ $$$$$$$$\ $$$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\  $$\   $$\ 
      $$ |$$ |\____$$  |\____$$  |$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\ $$ |  $$ |
$$\   $$ |$$ |  $$$$ _/   $$$$ _/ $$ |  $$ |$$$$$$$$ |$$ |  \__|$$ |  \__|$$ |  $$ |
$$ |  $$ |$$ | $$  _/    $$  _/   $$ |  $$ |$$   ____|$$ |      $$ |      $$ |  $$ |
\$$$$$$  |$$ |$$$$$$$$\ $$$$$$$$\ $$$$$$$  |\$$$$$$$\ $$ |      $$ |      \$$$$$$$ |
 \______/ \__|\________|\________|\_______/  \_______|\__|      \__|       \____$$ |
                                                                          $$\   $$ |
                                                                          \$$$$$$  |
                                                                           \______/`

func Init() error {
	err := initHelpers()
	if err != nil {
		return err
	}

	err = database.RunMigrations()
	if err != nil {
		return err
	}

	scrapers.RegisterScrapers()

	err = ffmpeg.IsExists()
	if err != nil {
		return err
	}

	err = IsFirstTime()
	if err != nil {
		return err
	}

	err = initWebApp()
	if err != nil {
		return err
	}
	return nil
}

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

func initWebApp() error {
	fmt.Println(art)
	router := mux.NewRouter()

	apps.RegisterFileServer(router)
	apps.RegisterApps(router)

	helpers.LogInfo("Server starting at "+addr, "Web")

	err := http.ListenAndServe(addr, router)
	if err != nil {
		return err
	}

	return nil
}

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
