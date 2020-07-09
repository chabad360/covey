package models

import (
	"database/sql/driver"
	"encoding/hex"
	"github.com/chabad360/covey/common"
	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"time"
)

type Map map[string]string

func (m Map) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *Map) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &m)
}

// Task defines the information of a task.
type Task struct {
	State     TaskState `json:"state" gorm:"<-:create;notnull"`
	Plugin    string    `json:"plugin" gorm:"<-:create;notnull"`
	ID        string    `json:"id" gorm:"<-:create;primarykey"`
	IDShort   string    `json:"-" gorm:"<-:create;notnull;unique"`
	NodeID    string    `json:"node" gorm:"<-:create;notnull"`
	Node      Node
	Details   Map       `json:"details" gorm:"<-:create;"`
	Log       []string  `json:"log" gorm:"<-:update;type:jsonb"`
	Time      time.Time `json:"time" gorm:"<-:create;notnull"`
	ExitCode  int       `json:"exit_code" gorm:"<-:create;notnull"`
	Command   string    `json:"-" gorm:"-"`
	JobID     string    `json:"job_id" gorm:"<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetIDShort returns the first 8 bytes of the task ID.
func (t *Task) GetIDShort() string { x, _ := hex.DecodeString(t.ID); return hex.EncodeToString(x[:8]) }

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.ExitCode = 258
	t.Log = []string{}
	t.State = StateQueued
	t.Time = time.Now()
	t.ID = common.GenerateID(t)
	t.IDShort = t.GetIDShort()
	return nil
}
