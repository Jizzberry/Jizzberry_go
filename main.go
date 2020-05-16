package main

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/logging"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	r *mux.Router
}

func main() {
	helpers.SetWorkingDirectory(".")
	helpers.CreateDirs()

	logging.Init()

	database.RunMigrations()
	scrapers.RegisterScrapers()
	err := ffmpeg.IsExists()
	if err != nil {
		return
	}

	router := mux.NewRouter()
	apps.RegisterApps(router)

	err = http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
	}

}
