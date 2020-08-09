package storage

import (
	"fmt"
	"github.com/chabad360/covey/models"
)

// TaskInfo contains new information about a running task.
type TaskInfo struct {
	Log      []string         `json:"log,omitempty"`
	ExitCode int              `json:"exit_code"`
	State    models.TaskState `json:"state"`
	ID       string           `json:"id"`
}

// AddTask adds a task to the database.
func AddTask(task *models.Task) error {
	return DB.Create(task).Error
}

// GetTask returns a task in the database.
func GetTask(id string) (*models.Task, bool) {
	var t models.Task
	if DB.Where("id = ?", id).Or("id_short = ?", id).First(&t).Error != nil {
		return nil, false
	}

	return &t, true
}

// SaveTask saves a live task to the database.
func SaveTask(t *TaskInfo) error {
	var task *models.Task
	var ok bool

	if task, ok = GetTask(t.ID); !ok {
		return fmt.Errorf("saveTask: task %s not found", t.ID)
	}

	if task.ExitCode == t.ExitCode && t.Log == nil && task.State == t.State { // Only update if there is something new!
		return fmt.Errorf("saveTask: nothing to save")
	}

	task.State = t.State
	task.ExitCode = t.ExitCode
	if t.Log != nil {
		task.Log = append(task.Log, t.Log...)
	}

	return DB.Save(task).Error
}
