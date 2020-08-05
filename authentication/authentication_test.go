package authentication

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

// revive:disable:line-length-limit
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIzIiwiYXVkIjoiYWxsIiwiZXhwIjo0NDY5NTQyMzYzLCJpYXQiOjE1OTU1MzI3NjMsImp0aSI6InNseUM0QVd2N0NhU3RKWG9yeXF5QzFPOWZLUFJzdFZQIn0.AAU9I8yub7VmTCnT833F54W6uQbhGVFKR8DSsi9pDJI"

// revive:enable:line-length-limit

func TestCreateToken(t *testing.T) {
	type args struct {
		userID        string
		tokenType     string
		allowedClaims []string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"user", args{"1", "user", []string{"all"}}, time.Now().Add(12 * time.Hour), false},
		{"api", args{"1", "api", []string{"all"}}, time.Now().Add(4 * (7 * (24 * time.Hour))), false},
		{"failBadID", args{"", "api", []string{"all"}}, time.Now().Add(4 * (7 * (24 * time.Hour))), true},
		{"failBadType", args{"1", "", []string{"all"}}, time.Now().Add(4 * (7 * (24 * time.Hour))), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := createToken(tt.args.userID, tt.args.tokenType, tt.args.allowedClaims)
			if (err != nil) != tt.wantErr {
				t.Errorf("createToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.Unix() != tt.want.Unix() {
				t.Errorf("createToken() got = %v, want %v", got.Unix(), tt.want.Unix())
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	losAngeles, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatal(err)
	}
	timeD := jwt.NumericDate(time.Date(2020, 7, 23, 12, 32, 43, 0, losAngeles))
	timeE := jwt.NumericDate(time.Date(2111, 8, 20, 11, 32, 43, 0, losAngeles))

	type args struct {
		tokenString string
		tokenType   string
	}
	tests := []struct {
		name    string
		args    args
		want    *jwt.Payload
		wantErr bool
	}{
		//revive:disable:line-length-limit
		{"Good", args{token, "api"}, &jwt.Payload{"covey-api", "3", []string{"all"}, timeE, nil, timeD, "slyC4AWv7CaStJXoryqyC1O9fKPRstVP"}, false},
		{"Expired", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIxIiwiYXVkIjoiYWxsIiwiZXhwIjowLCJpYXQiOjE1OTE5MTI5NzAsImp0aSI6InBIWWp4ZVVCclZmZHdVeldIUmloRkRQUkhCTXVFV21hIn0.XiNKXNDmsxXul8ceyQUgBWJBfrUmBsHWyLC34_Qy5qo",
			"api"}, &jwt.Payload{Issuer: "covey-api"}, true},
		{"Invalid", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIxIiwiYXVkIjoiYWxsIiwiZXhwIjowLCJpYXQiOjE1OTE5MTI5NzAsImp0aSI6InBIWWp4ZVVCclZmZHdVeldIUmloRkRQUkhCTXVFV21hIna.XiNKXNDmsxXul8ceyQUgBWJBfrUmBsHWyLC34_Qy5qo",
			"api"}, &jwt.Payload{Issuer: ""}, true},
		{"NoType", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjb3ZleS1hcGkiLCJzdWIiOiIxIiwiYXVkIjoiYWxsIiwiZXhwIjowLCJpYXQiOjE1OTE5MTI5NzAsImp0aSI6InBIWWp4ZVVCclZmZHdVeldIUmloRkRQUkhCTXVFV21hIna.XiNKXNDmsxXul8ceyQUgBWJBfrUmBsHWyLC34_Qy5qo",
			"asd"}, &jwt.Payload{Issuer: ""}, true},
		//revive:enable:line-length-limit
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseToken(tt.args.tokenString, tt.args.tokenType, "all")
			if (err != nil) != tt.wantErr {
				t.Errorf("parseToken() error = %v, wantErr %v, got = %v, sent = %v", err, tt.wantErr, got, tt.args)
			}
			if !tt.wantErr && !cmp.Equal(got, tt.want) {
				t.Errorf("parseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
