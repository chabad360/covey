package authentication

import (
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/chabad360/covey/test"
)

var (
	u = &models.User{
		Username: "user",
		Password: "password",
	}

	uu = &models.User{
		Username: "user",
		Password: "pass",
	}

	u2 = &models.User{
		Username: "user2",
		Password: "password",
	}
)

//revive:disable:cognitive-complexity
func TestTokenGetAPI(t *testing.T) {
	var tests = []struct {
		name       string
		userid     string
		want       string
		wantStatus int
	}{
		{"success", "1", `{"token":"`, http.StatusOK},
		{"fail", "", `{"error":"createToken: missing userID"}
`, http.StatusInternalServerError},
	}

	h := test.PureBoilerplate("GET", "/api/v1/auth/token", tokenGetAPI)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/auth/token", nil)
			if err != nil {
				t.Error(err)
			}
			if tt.userid != "" {
				req.Header.Add("X-User-ID", tt.userid)
			}

			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("tokenCookie status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}
			if !strings.Contains(rr.Body.String(), tt.want) {
				t.Errorf("tokenGetAPI body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	storage.DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	storage.AddUser(*uu)
	storage.AddUser(*u2)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
