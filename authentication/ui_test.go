package authentication

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/chabad360/covey/test"
)

//revive:disable:cognitive-complexity
func TestLogin(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(login).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("login status = %v, want %v", rr.Code, http.StatusInternalServerError)
	}
}

func TestLoginP(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		username   string
		password   string
		notWant    string
		url        string
		wantStatus int
	}{
		{"user", "pass", "", "/login", http.StatusFound},
		{"user", "pass", "", "/login?url=/home", http.StatusFound},
		{"us", "password", "a", "/login", http.StatusUnauthorized},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("POST", "/login", loginPost)

	for _, tt := range tests {
		testname := tt.username
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
				t.Errorf("loginP status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}

			if rr.Header().Get("set-cookie") == tt.notWant {
				t.Errorf("loginP cookie = %v, want: anything other than %v",
					rr.Header().Get("set-cookie"), tt.notWant)
			}
		})
	}
}

func Test_logout(t *testing.T) {
	h := test.PureBoilerplate("POST", "/logout", logout)

	rr, req, err := test.HTTPBoilerplate("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(rr, req)

	if rr.Header().Get("set-cookie") != "token=; Path=/; Expires=Sat, 01 Jan 2000 01:01:01 GMT; Max-Age=0" {
		t.Errorf("logout cookie = %v, want: %v", rr.Header().Get("set-cookie"),
			"token=; Path=/; Expires=Sat, 01 Jan 2000 01:01:01 GMT; Max-Age=0")
	}
}
