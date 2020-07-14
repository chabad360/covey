package models

import (
	"database/sql/driver"
	"encoding/hex"
	"github.com/chabad360/covey/common"
	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"time"
)

// TaskMap provides SQL scanner bindings for a map of JobTasks.
type TaskMap map[string]JobTask

// Value marshals a TaskMap to a SQL usable value
func (m TaskMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan unmarshals a stored value into a TaskMap
func (m *TaskMap) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &m)
}

// TaskArray provides SQL scanner bindings for an array of Tasks.
type TaskArray []Task

// Value marshals a TaskArray to a SQL usable value.
func (a TaskArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan unmarshals a stored value into a TaskArray.
func (a *TaskArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a)
}

// JobWithTasks is a regular Job, but with an expanded TaskHistory.
type JobWithTasks struct {
	Job
	TaskHistory TaskArray `json:"task_history"`
}

// JobTask represents a single task in a job.
type JobTask struct {
	Plugin  string    `json:"plugin"`
	Details StringMap `json:"details" gorm:"type:bytes"`
	Node    string    `json:"node,omitempty"`
}

// Job contains the information for a given job.
type Job struct {
	Name        string      `json:"name" gorm:"unique,notnull"`
	ID          string      `json:"id" gorm:"<-:create;primarykey"`
	IDShort     string      `json:"-" gorm:"<-:create;notnull;unique"`
	Cron        string      `json:"cron,omitempty"`
	Nodes       StringArray `json:"nodes" gorm:"notnull;type:bytes"`
	Tasks       TaskMap     `json:"tasks" gorm:"notnull;type:bytes"`
	TaskHistory StringArray `json:"task_history,omitempty" gorm:"<-:update;type:bytes"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// GetIDShort returns the first 8 bytes of the job ID.
func (j *Job) GetIDShort() string { x, _ := hex.DecodeString(j.ID); return hex.EncodeToString(x[:8]) }

// BeforeCreate generates the job ID before saving.
func (j *Job) BeforeCreate(tx *gorm.DB) (err error) {
	j.ID = common.GenerateID(j)
	j.IDShort = j.GetIDShort()
	return nil
}
