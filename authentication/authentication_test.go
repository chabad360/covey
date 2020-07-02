package authentication

import (
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

var token string

func TestCreateToken(t *testing.T) {
	type args struct {
		userid        string
		tokenType     string
		allowedClaims []string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"1", args{"1", "user", []string{"all"}}, time.Now().Add(20 * time.Minute), false},
		{"2", args{"1", "api", []string{"all"}}, time.Now().Add(4 * (7 * (24 * time.Hour))), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tk, got, err := createToken(tt.args.userid, tt.args.tokenType, tt.args.allowedClaims)
			if (err != nil) != tt.wantErr {
				t.Errorf("createToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Unix() != tt.want.Unix() {
				t.Errorf("createToken() got = %v, want %v", got.Unix(), tt.want.Unix())
			}
			token = tk
		})
	}
}

func TestParseToken(t *testing.T) {
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
		{"Good", args{token, "api"}, &jwt.Payload{Issuer: "covey-api"}, false},
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
				return
			}
			// if got != tt.want {
			// 	t.Errorf("parseToken() = %v, want %v", got, tt.want)
			// }
		})
	}
}
