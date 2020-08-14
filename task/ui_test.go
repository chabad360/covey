package task

import (
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
	"net/http"
	"testing"
)

func Test_uiTasks(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(uiTasks).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("uiTasks() status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func Test_uiTaskSingle(t *testing.T) {
	storage.DB.Delete(&models.Task{}, "id != ''")
	storage.AddTask(t2)

	//revive:disable:line-length-limit
	var tests = []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"success", "/tasks/" + t2.GetIDShort(), http.StatusOK},
		{"404", "/tasks/a83df6", http.StatusNotFound},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("GET", "/tasks/:task", uiTaskSingle)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("uiTaskSingle status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}
}

func Test_UITaskNew(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/new/task", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(UITaskNew).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("UITaskNew() status = %v, want %v", rr.Code, http.StatusOK)
	}
}
