package storage

import (
	"github.com/chabad360/covey/models"
)

// TaskInfo contains new information about a running task.
type TaskInfo struct {
	Log      []string `json:"log"`
	ExitCode int      `json:"exit_code"`
	ID       string   `json:"id"`
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
		return nil
	}

	if task.ExitCode == t.ExitCode || t.Log == nil { // Only update if there is something new!
		return nil
	}

	switch t.ExitCode {
	case 0:
		task.State = models.StateDone
	case 257:
		task.State = models.StateRunning
	default:
		task.State = models.StateError
	}

	task.ExitCode = t.ExitCode
	if t.Log != nil {
		task.Log = append(task.Log, t.Log...)
	}

	return DB.Save(task).Error
}