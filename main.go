package main

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks/rename"
	"github.com/gorilla/mux"
)

type App struct {
	r *mux.Router
}

func main() {
	//App:= App{r: mux.NewRouter()}

	helpers.SetWorkingDirectory("./")
	helpers.CreateDirs()
	database.RunMigrations()
	scrapers.RegisterScrapers()
	ffmpeg.IsExists()

	//fmt.Println(files.Initialize().Get(files.Files{GeneratedID: 1})[0].Tags)
	//scrapeActors.ScrapeActors{}.Start()
	//scan.Scan{}.Start(config.GetVideoPaths())
	//time.Sleep(10 * time.Minute)
	rename.Rename{}.Start(1, "just testing hello", []string{"Dillion Harper"})
}
