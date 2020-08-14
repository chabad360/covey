package node

import (
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
	"net/http"
	"testing"
)

func Test_uiNodes(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/nodes", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(uiNodes).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("uiNodes() status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func Test_uiNodeSingle(t *testing.T) {
	storage.DB.Delete(&models.Node{}, "id != ''")
	storage.AddNode(n2)

	//revive:disable:line-length-limit
	var tests = []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"success", "/nodes/n", http.StatusOK},
		{"404", "/nodes/node1", http.StatusNotFound},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("GET", "/nodes/:node", uiNodeSingle)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("uiNodeSingle status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}
}

func Test_UINodeNew(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/new/node", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(UINodeNew).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("UINodeNew() status = %v, want %v", rr.Code, http.StatusOK)
	}
}
