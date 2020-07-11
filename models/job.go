package models

import (
	"database/sql/driver"
	"encoding/hex"
	"github.com/chabad360/covey/common"
	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"time"
)

type TaskMap map[string]JobTask

func (m TaskMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}
func (m *TaskMap) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &m)
}

type TaskArray []Task

func (a TaskArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
func (a *TaskArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a)
}

type JobWithTasks struct {
	Job
	TaskHistory TaskArray `json:"task_history"`
}

// JobTask represents a single task in a job.
type JobTask struct {
	Plugin  string `json:"plugin"`
	Details Map    `json:"details" gorm:"type:bytes"`
	Node    string `json:"node,omitempty"`
}

// Job contains the information for a given job
type Job struct {
	Name        string    `json:"name" gorm:"unique,notnull"`
	ID          string    `json:"id" gorm:"<-:create;primarykey"`
	IDShort     string    `json:"-" gorm:"<-:create;notnull;unique"`
	Cron        string    `json:"cron,omitempty"`
	Nodes       Array     `json:"nodes" gorm:"notnull;type:bytes"`
	Tasks       TaskMap   `json:"tasks" gorm:"notnull;type:bytes"`
	TaskHistory Array     `json:"task_history" gorm:"<-:update;type:bytes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetIDShort returns the first 8 bytes of the job ID.
func (j *Job) GetIDShort() string { x, _ := hex.DecodeString(j.ID); return hex.EncodeToString(x[:8]) }

func (j *Job) BeforeCreate(tx *gorm.DB) (err error) {
	j.ID = common.GenerateID(j)
	j.IDShort = j.GetIDShort()
	return nil
}
