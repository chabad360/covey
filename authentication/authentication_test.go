package authentication

import (
	"testing"
	"time"
)

var token string

func Test_createToken(t *testing.T) {
	type args struct {
		userid        string
		tokenType     string
		allowedClaims map[string]bool
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"1", args{"1", "user", nil}, time.Now().Add(20 * time.Minute), false},
		{"2", args{"1", "api", nil}, time.Now().Add(4 * (7 * (24 * time.Hour))), false},
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

func Test_parseToken(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    *claims
		wantErr bool
	}{
		{"1", args{token}, &claims{Type: "api"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Type != tt.want.Type {
				t.Errorf("parseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
