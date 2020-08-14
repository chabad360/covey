package storage

import (
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/test"
	"github.com/google/go-cmp/cmp"
	"testing"
)

// revive:disable:line-length-limit
var n = test.N1

// revive:enable:line-length-limit

func TestAddNode(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want *models.Node
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", &n},
		{"3", &models.Node{}},
	}
	//revive:enable:line-length-limit

	testError := AddNode(&n)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Node
			if DB.Where("id = ?", tt.id).First(&got); cmp.Equal(got, tt.want) {
				t.Errorf("addNode() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestGetNodeIDOrName(t *testing.T) {
	DB.Delete(&models.Node{}, "id != ''")
	AddNode(&n)
	var tests = []struct {
		name  string
		id    string
		field string
		want  string
		want2 bool
	}{
		{"ok_ID", n.Name, "id", n.ID, true},
		{"notok_ID", "n", "id", "", false},
		{"ok_Name", n.ID, "name", n.Name, true},
		{"notok_Name", "n", "name", "", false},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			got, got2 := GetNodeIDorName(tt.id, tt.field)
			if got2 != tt.want2 {
				t.Errorf("GetNodeIDorName() = %v, want %v", got2, tt.want2)
			}
			if got != tt.want {
				t.Errorf("GetNodeIDorName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNode(t *testing.T) {
	DB.Delete(&models.Node{}, "id != ''")
	AddNode(&n)
	type args struct {
		id string
	}
	tests := []struct {
		name  string
		args  args
		want  *models.Node
		want1 bool
	}{
		{"success", args{n.ID}, &n, true},
		{"fail", args{"3"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := GetNode(tt.args.id)
			if got1 != tt.want1 {
				t.Errorf("GetNode() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
