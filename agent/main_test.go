package main

import "testing"

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := newRunningTask(task{"echo hello world", "test"})
		run(t)
	}
}
