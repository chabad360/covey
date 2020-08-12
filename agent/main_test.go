package main

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := newRunningTask(task{"echo hello world", "test"})
		run(t)
	}
}

func Test_genBody(t *testing.T) {
	t1 := newRunningTask(task{"test", "test1"})
	t2 := newRunningTask(task{"test", "test2"})
	t2.Log("test")
	t3 := newRunningTask(task{"test", "test3"})
	t3.Finish(0, 0)

	type args struct {
		rt *runningTask
	}
	tests := []struct {
		name    string
		args    args
		want    *runningTask
		want1   []byte
		wantErr bool
	}{
		{"new", args{t1}, t1, []byte(`{"exit_code":257,"state":2,"id":"test1"}`), false},
		{"Log", args{t2}, t2, []byte(`{"log":["test"],"exit_code":257,"state":2,"id":"test2"}`), false},
		{"Done", args{t3}, nil, []byte(`{"exit_code":0,"state":0,"id":"test3"}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := genBody(tt.args.rt)
			if (err != nil) != tt.wantErr {
				t.Errorf("genBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genBody() got = %v, want %v", got, tt.want)
			}
			if !cmp.Equal(got1, tt.want1) {
				t.Errorf("genBody() got1 = %v, want %v", string(got1), string(tt.want1))
			}
		})
	}
}

func Test_run(t *testing.T) {
	t1 := newRunningTask(task{"echo hello", "test1"})
	t2 := newRunningTask(task{"echo 'hello\nworld'", "test2"})
	t3 := newRunningTask(task{"exit 1'", "test"})

	type args struct {
		t *runningTask
	}
	tests := []struct {
		name         string
		args         args
		wantLog      []string
		wantExitCode int
		wantState    int
	}{
		{"success", args{t1}, []string{"hello"}, 0, 0},
		{"successNewLine", args{t2}, []string{"hello", "world"}, 0, 0},
		{"successExitError", args{t3}, nil, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run(tt.args.t)
			<-tt.args.t.Done()

			gotLog := tt.args.t.GetLog()
			gotExitCode := tt.args.t.ExitCode
			gotState := tt.args.t.State
			if !cmp.Equal(gotLog, tt.wantLog) {
				t.Errorf("run() gotLog = %v, want %v", gotLog, tt.wantLog)
			}
			if !cmp.Equal(gotExitCode, tt.wantExitCode) {
				t.Errorf("run() gotExitCode = %v, want %v", gotExitCode, tt.wantExitCode)
			}
			if !cmp.Equal(gotState, tt.wantState) {
				t.Errorf("run() gotState = %v, want %v", gotState, tt.wantState)
			}
		})
	}
}
