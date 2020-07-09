package models

import (
	"encoding/hex"
	"github.com/chabad360/covey/common"
	"gorm.io/gorm"
)

// JobTask represents a single task in a job.
type JobTask struct {
	Plugin  string            `json:"plugin"`
	Details map[string]string `json:"details"`
	Node    string            `json:"node,omitempty"`
}

// Job contains the information for a given job
type Job struct {
	Name        string             `json:"name" gorm:"unique,notnull"`
	ID          string             `json:"id" gorm:"<-:create;primarykey"`
	IDShort     string             `json:"-" gorm:"<-:create;notnull;unique"`
	Cron        string             `json:"cron,omitempty"`
	Nodes       []string           `json:"nodes" gorm:"notnull"`
	Tasks       map[string]JobTask `json:"tasks"`
	TaskHistory []Task             `json:"task_history" gorm:"<-:update"`
}

// GetIDShort returns the first 8 bytes of the job ID.
func (j *Job) GetIDShort() string { x, _ := hex.DecodeString(j.ID); return hex.EncodeToString(x[:8]) }

func (j *Job) BeforeCreate(tx *gorm.DB) (err error) {
	j.ID = common.GenerateID(j)
	j.IDShort = j.GetIDShort()
	return nil
}
