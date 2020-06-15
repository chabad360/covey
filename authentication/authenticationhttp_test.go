package authentication

import (
	"fmt"
	"net/http"
	"strings"
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
		url        string
		wantStatus int
	}{
		{"user", "password", "", "/auth/login", http.StatusFound},
		{"user", "password", "", "/auth/login?url=/home", http.StatusFound},
		{"us", "password", "a", "/auth/login", http.StatusUnauthorized},
		{"", "", "a", "/auth/login", http.StatusForbidden},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("POST", "/auth/login", tokenGet)

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.username)
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("POST", tt.url, strings.NewReader(
				fmt.Sprintf("username=%s&password=%s", tt.username, tt.password)))
			if tt.username == "" {
				rr, req, err = test.HTTPBoilerplate("POST", tt.url, nil)
			}
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	if rr.Header().Get("set-cookie") != "token=; Path=/; Expires=Sat, 01 Jan 2000 01:01:01 GMT; Max-Age=0" {
		t.Errorf("tokenRemove cookie = %v, want: %v", rr.Header().Get("set-cookie"),
			"token=; Path=/; Expires=Sat, 01 Jan 2000 01:01:01 GMT; Max-Age=0")
	}
}
