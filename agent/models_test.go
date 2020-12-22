package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_queue_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    task
		wantErr bool
	}{
		{"success", args{[]byte(`{"0":{"command":"test","id":"test1"}}`)}, task{"test", "test1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q = &Queue{}

			if err := q.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := q.GetNext(); !cmp.Equal(got, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
