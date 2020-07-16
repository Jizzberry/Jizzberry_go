package manager

import (
	"context"
	"fmt"
)

type Task struct {
	Uid      string              `json:"uid"`
	Name     string              `json:"name"`
	Cancel   *context.CancelFunc `json:"-"`
	Progress *int                `json:"progress"`
}

func isTaskActive(uid string) bool {
	if task, ok := GlobalTasksStorage[uid]; ok {
		if *task.Progress != 100 {
			return true
		}
	}
	return false
}

func GetProgress(uid string) int {
	if task, ok := GlobalTasksStorage[uid]; ok {
		if *task.Progress > 99 {
			removeTask(uid)
		}
		return *task.Progress
	}
	return -1
}

func GetAllTaskStatus() (tasks []Task) {
	for key, value := range GlobalTasksStorage {
		tasks = append(tasks, *value)
		if *value.Progress > 99 {
			removeTask(key)
		}
	}
	return
}

func StopTask(uid string) error {
	if val, ok := GlobalTasksStorage[uid]; ok {
		cancel := *val.Cancel
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
	delete(GlobalTasksStorage, uid)
}
