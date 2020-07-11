package models

import (
	"testing"
)

var ta = &Task{
	ID:       "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1",
	Plugin:   "plugin",
	State:    StateDone,
	Details:  map[string]string{"test": "test"},
	Log:      []string{"hello", "world"},
	ExitCode: 1,
}

func TestTask_GetIDShort(t *testing.T) {
	if got := ta.GetIDShort(); got != "a7a39b72f29718e6" {
		t.Errorf("Task.GetIDShort() = %v, want %v", got, "a7a39b72f29718e6")
	}
}
