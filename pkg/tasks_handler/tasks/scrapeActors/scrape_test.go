package scrapeActors_test

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks/scrapeActors"
	"io/ioutil"
	"testing"
	"time"
)

func TestScrapeActors(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")

	t.Log(dir)

	helpers.SetWorkingDirectory(dir)
	helpers.CreateDirs()
	scrapers.RegisterScrapers()

	database.RunMigrations()

	task := scrapeActors.ScrapeActors{}
	cancel, _ := task.Start()
	tmp := *cancel
	time.Sleep(10 * time.Second)

	defer tmp()
}
