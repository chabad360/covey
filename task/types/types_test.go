package types

import (
	"reflect"
	"testing"
	"time"
)

var n = &Task{
	ID:       "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1",
	Plugin:   "plugin",
	State:    StateDone,
	Node:     "node",
	Time:     time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
	Details:  map[string]string{"test": "test"},
	Log:      []string{"hello", "world"},
	ExitCode: 1,
}

func TestTask_GetID(t *testing.T) {
	if got := n.GetID(); got != "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1" {
		t.Errorf("Task.GetID() = %v, want %v", got, "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1")
	}
}

func TestTask_GetIDShort(t *testing.T) {
	if got := n.GetIDShort(); got != "a7a39b72f29718e6" {
		t.Errorf("Task.GetIDShort() = %v, want %v", got, "a7a39b72f29718e6")
	}
}

func TestTask_GetPlugin(t *testing.T) {
	if got := n.GetPlugin(); got != "plugin" {
		t.Errorf("Task.GetPlugin() = %v, want %v", got, "plugin")
	}
}

func TestTask_GetState(t *testing.T) {
	if got := n.GetState(); got != n.State {
		t.Errorf("Task.GetState() = %v, want %v", got, StateDone)
	}
}

func TestTask_GetNode(t *testing.T) {
	if got := n.GetNode(); got != "node" {
		t.Errorf("Task.GetNode() = %v, want %v", got, "node")
	}
}

func TestTask_GetDetails(t *testing.T) {
	if got := n.GetDetails(); !reflect.DeepEqual(got, n.Details) {
		t.Errorf("Task.GetDetails() = %v, want %v", got, n.Details)
	}
}

func TestTask_GetTime(t *testing.T) {
	if got := n.GetTime(); got != n.Time {
		t.Errorf("Task.GetTime() = %v, want %v", got, n.Time)
	}
}

func TestTask_GetExitCode(t *testing.T) {
	if got := n.GetExitCode(); got != n.ExitCode {
		t.Errorf("Task.GetExitCode() = %v, want %v", got, n.ExitCode)
	}
}

func TestTask_GetLog(t *testing.T) {
	if got := n.GetLog(); got[1] != "world" {
		t.Errorf("Task.GetLog() = %v, want %v", got, []string{"hello", "world"})
	}
}
