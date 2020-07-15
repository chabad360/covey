package job

import (
	"testing"

	"github.com/chabad360/covey/storage"
)

func TestGetJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want bool
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", true},
		{"3", false},
	}
	//revive:enable:line-length-limit
	storage.DB = db

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if _, ok := GetJob(tt.id); ok != tt.want {
				t.Errorf("GetJob() = %v, want %v", ok, tt.want)
			}
		})
	}
}

func TestGetJobWithTasks(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want bool
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", true},
		{"3", false},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if _, ok := GetJob(tt.id); ok != tt.want {
				t.Errorf("GetJobWithTasks() = %v, want %v", ok, tt.want)
			}
		})
	}
}
