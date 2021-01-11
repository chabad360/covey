package authentication

import (
	"net/http"
	"testing"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/test"
)

//revive:disable:cognitive-complexity

func TestAuthUserMiddleware(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		name       string
		url        string
		token      string
		location   string
		wantStatus int
	}{
		{"redirect", "/test", "", "/login?url=/test", http.StatusFound},
		{"success", "/test", test.JWT5, "", http.StatusOK},
		{"ignored", "/login", test.JWT5, "", http.StatusOK},
		{"logout", "/test", "1", "/logout", http.StatusFound},
		{"logout2", "/logout", "1", "", http.StatusOK},
		{"login", "/login", "", "", http.StatusOK},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.name
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

func TestAuthAPIMiddleware(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		name       string
		token      string
		want       string
		wantStatus int
	}{
		{"forbidden", "", `{"error":"invalid bearer token"}
`, http.StatusBadRequest},
		{"fail", "123", `{"error":"jwt: malformed token"}
`, http.StatusUnauthorized},
		{"success", test.JWT1, "3", http.StatusOK},
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
