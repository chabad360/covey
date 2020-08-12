package main

import (
	"os"
	"testing"
)

func Test_settings(t *testing.T) {
	m1 := make(map[string]string)
	m1["AGENT_ID"] = "test"
	m1["AGENT_HOST"] = "127.0.0.1"

	type args struct {
		conf *config
	}
	tests := []struct {
		name    string
		args    args
		env     map[string]string
		aPath   string
		wantErr bool
	}{
		{"success", args{&config{}}, m1, "http://127.0.0.1:8080/agent/test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			if err := settings(tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("settings() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.conf.AgentPath != tt.aPath {
				t.Errorf("settings() agent path = %v, want %v", tt.args.conf.AgentPath, tt.aPath)
			}
		})
	}
}
