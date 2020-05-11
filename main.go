package main

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
)

func main() {
	helpers.SetWorkingDirectory("./")
	helpers.CreateDirs()
	database.RunMigrations()
	scrapers.RegisterScrapers()
	ffmpeg.IsExists()
}
