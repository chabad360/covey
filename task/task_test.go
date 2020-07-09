package task

import (
	"github.com/chabad360/covey/models"
	"testing"
)

// var task1 = &types.Task{
// 	ID:       "2778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6c",
// 	State:    types.StateRunning,
// 	Plugin:   "test",
// 	Node:     "test",
// 	Time:     time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
// 	Details:  map[string]string{"test": "test"},
// 	ExitCode: 0,
// }

func TestGetTask(t *testing.T) {
	type args struct {
		identifier string
	}
	tests := []struct {
		name  string
		args  args
		want  *models.Task
		want1 bool
	}{
		{"db", args{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e"}, task, true},
		{"noDB", args{"31b079725d0a20bfe6c3b6e"}, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getTask(tt.args.identifier)
			if (got1 != nil) == tt.want1 && got.ID != tt.want.ID {
				t.Errorf("getTask() got = %v, want %v, got1 = %v, want %v", got, tt.want, got1, tt.want1)
			}
		})
	}
}
