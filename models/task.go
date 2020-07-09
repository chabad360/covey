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
	return json.Marshal(m)
}
func (m *Map) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &m)
}

type Array []string

func (a Array) Value() (driver.Value, error) {
	return json.Marshal(a)
}
func (a *Array) Scan(value interface{}) error {
	//b, ok := value.([]byte)
	//if !ok {
	//	return fmt.Errorf("expected []byte, got %v", value)
	//}
	return json.Unmarshal(value.([]byte), &a)
}

// Task defines the information of a task.
type Task struct {
	State     TaskState `json:"state" gorm:"notnull"`
	Plugin    string    `json:"plugin" gorm:"<-:create;notnull"`
	ID        string    `json:"id" gorm:"<-:create;primarykey"`
	IDShort   string    `json:"-" gorm:"<-:create;notnull;unique"`
	NodeID    string    `json:"node" gorm:"<-:create;notnull"`
	Node      Node      `json:"-" gorm:"<-:create;"`
	Details   Map       `json:"details" gorm:"<-:create;"`
	Log       Array     `json:"log" gorm:"<-:update"`
	Time      time.Time `json:"time" gorm:"<-:create;notnull"`
	ExitCode  int       `json:"exit_code" gorm:"notnull"`
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
