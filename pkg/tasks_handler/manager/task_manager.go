package manager

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/config"
	"github.com/Jizzberry/Jizzberry-go/pkg/logging"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks/rename"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks/scan"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks/scrapeActors"
	"github.com/google/uuid"
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
	uid, err := uuid.NewRandom()
	if err != nil {
		logging.LogError(err.Error(), "Manager - Scan")
	}
	cancel, progress := scan.Scan{}.Start(config.GetVideoPaths())
	GlobalTasksStorage.Cancel[uid.String()] = cancel
	GlobalTasksStorage.Progress[uid.String()] = progress
	return uid.String()
}

func StartScrapeActors() string {
	uid, err := uuid.NewRandom()
	if err != nil {
		logging.LogError(err.Error(), "Manager - ScrapeActors")
	}
	cancel, progress := scrapeActors.ScrapeActors{}.Start()
	GlobalTasksStorage.Cancel[uid.String()] = cancel
	GlobalTasksStorage.Progress[uid.String()] = progress
	return uid.String()
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
	return 0
}

func StopTask(uid string) {
	if val, ok := GlobalTasksStorage.Cancel[uid]; ok {
		cancel := *val
		cancel()
	}
}
