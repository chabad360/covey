package models

import (
	"database/sql/driver"
	"encoding/hex"
	"time"

	json "github.com/json-iterator/go"
	"gorm.io/gorm"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/models/safe"
)

// StringMap provides SQL scanner bindings for a map of strings.
type StringMap map[string]string

// Value marshals a StringMap to a SQL usable value
func (m StringMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan unmarshals a stored value into a StringMap
func (m *StringMap) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &m)
}

// StringArray provides SQL scanner bindings for an array of strings.
type StringArray []string

// Value marshals a StringArray to a SQL usable value
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan un-marshals a stored value into a StringArray
func (a *StringArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a)
}

// Task defines the information of a task.
type Task struct {
	State     TaskState   `json:"state" gorm:"notnull"`
	Plugin    string      `json:"plugin" gorm:"<-:create;notnull"`
	ID        string      `json:"id" gorm:"<-:create;primarykey"`
	IDShort   string      `json:"-" gorm:"<-:create;notnull;unique"`
	Node      string      `json:"node" gorm:"<-:create;notnull"`
	Details   StringMap   `json:"details" gorm:"<-:create;"`
	Log       StringArray `json:"log,omitempty" gorm:"<-:update;type:bytea"`
	ExitCode  int         `json:"exit_code" gorm:"notnull"`
	Command   string      `json:"-" gorm:"-"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// GetIDShort returns the first 8 bytes of the task ID.
func (t *Task) GetIDShort() string { x, _ := hex.DecodeString(t.ID); return hex.EncodeToString(x[:8]) }

// BeforeCreate initializes the default values for a Task.
func (t *Task) BeforeCreate(_ *gorm.DB) (err error) {
	t.ExitCode = 258
	t.State = StateQueued
	t.ID = common.GenerateID(t)
	t.IDShort = t.GetIDShort()
	return nil
}

// ToSafe converts a Task into a plugin-safe safe.Task wrapper.
func (t *Task) ToSafe() safe.Task {
	return safe.Task{
		Plugin:  t.Plugin,
		Node:    t.Node,
		Details: t.Details,
	}
}
