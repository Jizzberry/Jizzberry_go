package main

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	r *mux.Router
}

func main() {
	helpers.SetWorkingDirectory("./")
	helpers.CreateDirs()
	database.RunMigrations()
	scrapers.RegisterScrapers()
	ffmpeg.IsExists()

	//fmt.Println(actor_details.Initialize().Get(actor_details.ActorDetails{}))
	//scrapeActors.ScrapeActors{}.Start()
	//scan.Scan{}.Start(config.GetVideoPaths())

	//uid := manager.StartScan()
	//for i := 0; i < 100; i++ {
	//	manager.GetProgress(uid)
	//	time.Sleep(2 * time.Second)
	//}

	//rename.Rename{}.Start(1, "just testing hello", []string{"Dillion Harper"})

	router := mux.NewRouter()
	apps.RegisterApps(router)

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
	}

	//auth.Initialize().Create(auth.Auth{
	//	Username: "ovenoboyo",
	//	Password: "kekboi69",
	//})
	//
	//fmt.Println(auth.Initialize().Get(auth.Auth{Username: "ovenoboyo"}))
}
