package initializer

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/apps"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/gorilla/mux"
	"net/http"
)

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

	err = initWebApp()
	if err != nil {
		return err
	}
	return nil
}

func initHelpers() error {
	err := helpers.CreateDirs()
	if err != nil {
		return err
	}
	err = helpers.ConfigInit()
	if err != nil {
		return err
	}
	helpers.LoggerInit()
	helpers.RndInit()
	return nil
}

func initWebApp() error {
	router := mux.NewRouter()

	apps.RegisterFileServer(router)
	apps.RegisterApps(router)

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		return err
	}
	return nil
}
