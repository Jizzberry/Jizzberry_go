package manager

import (
	"context"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/tasks/rename"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/tasks/scan"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/tasks/scrapeActors"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/tasks/scrapeStudios"
	"strconv"
	"time"
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

	cancel, progress := scan.Scan{}.Start(helpers.GetConfig().Paths)
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

func StartScrapeStudios() string {
	uid := "scrapeStudiosList"

	if isTaskActive(uid) {
		return uid
	}

	cancel, progress := scrapeStudios.ScrapeStudios{}.Start()
	GlobalTasksStorage.Cancel[uid] = cancel
	GlobalTasksStorage.Progress[uid] = progress
	return uid
}

//Rename shouldn't be stopped
func StartRename(sceneId int64) string {
	uid := strconv.FormatInt(time.Now().Unix(), 10)
	cancel, progress := rename.Rename{}.Start(sceneId)
	GlobalTasksStorage.Cancel[uid] = cancel
	GlobalTasksStorage.Progress[uid] = progress
	return uid

}

func GetProgress(uid string) int {
	if val, ok := GlobalTasksStorage.Progress[uid]; ok {
		progress := *val
		if progress > 99 {
			removeTask(uid)
		}
		return progress
	}
	return -1
}

func GetAllProgress() map[string]int {
	tmpMap := make(map[string]int)
	for key, value := range GlobalTasksStorage.Progress {
		tmpMap[key] = *value
		if *value > 99 {
			removeTask(key)
		}
	}
	return tmpMap
}

func StopTask(uid string) error {
	if val, ok := GlobalTasksStorage.Cancel[uid]; ok {
		cancel := *val
		if cancel != nil {
			cancel()
			removeTask(uid)
			return nil
		}
		return fmt.Errorf("task not stoppable")
	}
	return fmt.Errorf("task does not exist")
}

func removeTask(uid string) {
	delete(GlobalTasksStorage.Progress, uid)
	delete(GlobalTasksStorage.Cancel, uid)
}
