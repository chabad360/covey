package authentication

import (
	"net/http"
	"strings"
	"testing"

	"github.com/chabad360/covey/test"
)

//revive:disable:cognitive-complexity
func TestTokenGetAPI(t *testing.T) {
	var tests = []struct {
		userid     string
		want       string
		wantStatus int
	}{
		{"1", `{"token":"`, http.StatusOK},
		{"", `{"error":"createToken: missing userID"}
`, http.StatusInternalServerError},
	}

	h := test.PureBoilerplate("GET", "/api/v1/auth/token", tokenGetAPI)

	for _, tt := range tests {
		testname := tt.userid
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/auth/token", nil)
			if err != nil {
				t.Fatal(err)
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
