package manager

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	files2 "github.com/Jizzberry/Jizzberry_go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/tasks"
	"strconv"
	"time"
)

var GlobalTasksStorage = make(map[string]*Task)

func StartScan() string {
	uid := "scan"

	if isTaskActive(uid) {
		return uid
	}

	cancel, progress := tasks.Scan{}.Start(helpers.GetConfig().Paths)
	GlobalTasksStorage[uid] = &Task{
		Uid:      uid,
		Name:     "Scan",
		Cancel:   cancel,
		Progress: progress,
	}
	return uid
}

func StartScrapeActors() string {
	uid := "scrapeActorsList"

	if isTaskActive(uid) {
		return uid
	}

	cancel, progress := tasks.ScrapeActors{}.Start()
	GlobalTasksStorage[uid] = &Task{
		Uid:      uid,
		Name:     "Scrape Actors List",
		Cancel:   cancel,
		Progress: progress,
	}
	return uid
}

func StartScrapeStudios() string {
	uid := "scrapeStudiosList"

	if isTaskActive(uid) {
		return uid
	}

	cancel, progress := tasks.ScrapeStudios{}.Start()
	GlobalTasksStorage[uid] = &Task{
		Uid:      uid,
		Name:     "Scrape Studios List",
		Cancel:   cancel,
		Progress: progress,
	}
	return uid
}

//Rename shouldn't be stopped
func StartRename(sceneId int64) string {
	uid := strconv.FormatInt(time.Now().Unix(), 10)
	cancel, progress := tasks.Rename{}.Start(sceneId)
	GlobalTasksStorage[uid] = &Task{
		Uid:      uid,
		Name:     "Rename",
		Cancel:   cancel,
		Progress: progress,
	}
	return uid
}

func OrganiseAll() {
	model := files2.Initialize()
	defer model.Close()

	file := model.Get(files2.Files{})
	for _, f := range file {
		StartRename(f.SceneID)
	}
}
