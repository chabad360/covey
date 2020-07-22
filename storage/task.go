package storage

import (
	"errors"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"
)

// TaskInfo contains new information about a running task.
type TaskInfo struct {
	Log      []string `json:"log"`
	ExitCode int      `json:"exit_code"`
	ID       string   `json:"id"`
}

// AddTask adds a task to the database.
func AddTask(task *models.Task) error {
	result := DB.Create(task)
	return result.Error
}

// GetTask returns a task in the database.
func GetTask(id string) (*models.Task, bool) {
	var t models.Task
	result := DB.Where("id = ?", id).Or("id_short = ?", id).First(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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

	result := DB.Save(task)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return nil
}
