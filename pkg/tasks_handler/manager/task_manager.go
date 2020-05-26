package manager

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks/rename"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks/scan"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks/scrapeActors"
)

type TasksStorage struct {
	Cancel   map[string]*context.CancelFunc
	Progress map[string]*int
}

var GlobalTasksStorage = TasksStorage{
	Cancel:   make(map[string]*context.CancelFunc),
	Progress: make(map[string]*int),
}

func StartScan() string {
	uid := "scan"

	if isTaskActive(uid) {
		return uid
	}

	cancel, progress := scan.Scan{}.Start(helpers.GetVideoPaths())
	GlobalTasksStorage.Cancel[uid] = cancel
	GlobalTasksStorage.Progress[uid] = progress
	return uid
}

func isTaskActive(uid string) bool {
	if _, ok := GlobalTasksStorage.Cancel[uid]; ok {
		if val, ok := GlobalTasksStorage.Progress[uid]; ok {
			if *val != 100 {
				return true
			}
		}
	}
	return false
}

func StartScrapeActors() string {
	uid := "scrapeActorsList"

	if isTaskActive(uid) {
		return uid
	}

	cancel, progress := scrapeActors.ScrapeActors{}.Start()
	GlobalTasksStorage.Cancel[uid] = cancel
	GlobalTasksStorage.Progress[uid] = progress
	return uid
}

//Rename shouldn't be stopped
func StartRename(sceneId int64, title string, actors []string) {
	rename.Rename{}.Start(sceneId, title, actors)
}

func GetProgress(uid string) int {
	if val, ok := GlobalTasksStorage.Progress[uid]; ok {
		progress := *val
		return progress
	}
	return -1
}

func StopTask(uid string) {
	if val, ok := GlobalTasksStorage.Cancel[uid]; ok {
		cancel := *val
		cancel()
	}
}
