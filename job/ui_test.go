package job

import (
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
	"net/http"
	"testing"
)

func Test_uiJobs(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/jobs", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(uiJobs).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("uiJobs() status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func Test_uiJobSingle(t *testing.T) {
	storage.DB.Delete(&models.Job{}, "id != ''")
	storage.AddJob(&j2)

	//revive:disable:line-length-limit
	var tests = []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"success", "/jobs/add", http.StatusOK},
		{"404", "/jobs/update", http.StatusNotFound},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("GET", "/jobs/:job", uiJobSingle)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("uiJobSingle status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}
}

func Test_UIJobNew(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/new/job", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(UIJobNew).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("UIJobNew() status = %v, want %v", rr.Code, http.StatusOK)
	}
}
