package models

// TaskState represents the current state of the task.
type TaskState int

const (
	// StateError should be given if the task fails.
	StateError TaskState = iota
	// StateDone should be given if the task is complete.
	StateDone
	// StateRunning should be given if the task is running.
	StateRunning
	// StateAborted should be given if the task is aborted (stopped while running).
	StateAborted
	// StateCancelled should be given if the queued task is canceled.
	StateCancelled
	// StateScheduled should be given if the task is scheduled to be executed.
	// StateScheduled = 5
	_
	// StateQueued should be given while the task is waiting to be executed.
	StateQueued
	// StateOther should be avoided.
	StateOther TaskState = 10
	// StateInternalError should only be given if the task plugin fails.
	// This will be used if err != nil, if this is given there is a bug in the plugin.
	StateInternalError TaskState = 11
)
