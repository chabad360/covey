package models

import "testing"

var j = &Job{
	Name: "test",
	ID:   "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1",
	// Cron:        "0 * * * * *",
	// Nodes:       tt.fields.Nodes,
	// Tasks:       tt.fields.Tasks,
	// TaskHistory: tt.fields.TaskHistory,
}

func TestJob_GetIDShort(t *testing.T) {
	if got := j.GetIDShort(); got != "a7a39b72f29718e6" {
		t.Errorf("Job.GetIDShort() = %v, want %v", got, "a7a39b72f29718e6")
	}
}
