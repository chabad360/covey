package authentication

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/chabad360/covey/test"
)

//revive:disable:cognitive-complexity
func Test_tokenGet(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		username   string
		password   string
		notWant    string
		wantStatus int
	}{
		{"user", "password", "", http.StatusFound},
		{"us", "password", "a", http.StatusUnauthorized},
		{"", "", "a", http.StatusForbidden},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("GET", "/auth/login", tokenGet)

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.username)
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/auth/login", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.username != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}

			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("tokenGet status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}

			if rr.Header().Get("set-cookie") == tt.notWant {
				t.Errorf("tokenGet cookie = %v, want: anything other than %v",
					rr.Header().Get("set-cookie"), tt.notWant)
			}
		})
	}
}

func Test_tokenRemove(t *testing.T) {
	h := test.PureBoilerplate("POST", "/auth/logout", tokenRemove)

	rr, req, err := test.HTTPBoilerplate("POST", "/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(rr, req)

	if rr.Header().Get("set-cookie") != "token=; Max-Age=0" {
		t.Errorf("tokenRemove cookie = %v, want: %v", rr.Header().Get("set-cookie"), "token=; Max-Age=0")
	}
}
