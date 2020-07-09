package task

import (
	"errors"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"

	"github.com/chabad360/covey/storage"
)

var db *gorm.DB

// AddTask adds a task to the database.
func addTask(task *models.Task) error {
	refreshDB()

	result := db.Create(task)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return nil
}

// GetTask returns a task in the database.
func getTask(id string) (*models.Task, bool) {
	refreshDB()

	var t models.Task
	result := db.Where("id = ?", id).Or("id_short = ?", id).First(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false
	}

	return &t, true
}

// saveTask saves a live task to the database.
func saveTask(t *TaskInfo) error {
	refreshDB()

	var task *models.Task
	var ok bool

	if task, ok = getTask(t.ID); !ok {
		return nil
	}

	if task.ExitCode != t.ExitCode || t.Log != nil { // Only update if there is something new!
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

		result := db.Save(task)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
	}

	return nil
}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
