package main

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/auth"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	helpers.SetWorkingDirectory(".")
	helpers.Init()

	database.RunMigrations()
	scrapers.RegisterScrapers()
	err := ffmpeg.IsExists()
	if err != nil {
		return
	}

	router := mux.NewRouter()
	apps.RegisterFileServer(router)
	apps.RegisterApps(router)

	auth.Initialize().Create(auth.Auth{
		Username: "ovenoboyo",
		Password: "kekboi69",
		IsAdmin:  true,
	})

	err = http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
	}
}
