package authentication

import (
	"net/http"
	"testing"

	"github.com/chabad360/covey/test"
	"github.com/go-playground/pure/v5"
)

var tokenUser, _, _ = createToken("1", "user", []string{"all"})

//revive:disable:cognitive-complexity

func Test_AuthUserMiddlware(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		url        string
		token      string
		location   string
		wantStatus int
	}{
		{"/test", "", "/login?url=/test", http.StatusFound},
		{"/test", tokenUser, "", http.StatusOK},
		{"/login", tokenUser, "/dashboard", http.StatusFound},
		{"/test", "1", "/logout", http.StatusFound},
		{"/logout", "1", "", http.StatusOK},
		{"/login", "", "", http.StatusOK},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.token
		t.Run(testname, func(t *testing.T) {
			p := pure.New()
			p.Use(AuthUserMiddleware)
			p.Get(tt.url, func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("success")) })
			h := p.Serve()

			rr, req, err := test.HTTPBoilerplate("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.token != "" {
				req.Header.Add("cookie", "token="+tt.token)
			}

			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("AuthUserMiddlware status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}

			if rr.Header().Get("location") != tt.location {
				t.Errorf("AuthUserMiddlware location = %v, want %v", rr.Header().Get("location"), tt.location)
			}
		})
	}
}

func Test_AuthAPIMiddlware(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		name       string
		token      string
		want       string
		wantStatus int
	}{
		{"forbidden", "", `{"error":"forbidden"}
`, http.StatusForbidden},
		{"fail", "123", `{"error":"jwt: malformed token"}
`, http.StatusUnauthorized},
		{"success", tokenUser, "1", http.StatusOK},
	}
	//revive:enable:line-length-limit
	p := pure.New()
	p.Use(AuthAPIMiddleware)
	p.Get("/test", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(r.Header.Get("x-user-id"))) })
	h := p.Serve()

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.token != "" {
				req.Header.Add("Authorization", "Bearer "+tt.token)
			}

			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("AuthAPIMiddlware status = %v, want %v, error = %v", rr.Code, tt.wantStatus, rr.Body.String())
			}

			if rr.Body.String() != tt.want {
				t.Errorf("AuthAPIMiddlware body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}
