package models

// TaskState represents the current state of the task.
type TaskState int

const (
	// StateDone should be given if the task is complete.
	StateDone TaskState = iota
	// StateError should be given if the task fails.
	StateError
	// StateRunning should be given if the task is running.
	StateRunning
	// StateAborted should be given if the task is aborted (stopped while running).
	StateAborted
	// StateCancelled should be given if the queued task is canceled.
	StateCancelled
	// StateScheduled should be given if the task is scheduled to be executed.
	// StateScheduled
	_
	// StateQueued should be given while the task is waiting to be executed.
	StateQueued
	// StateOther should be avoided.
	StateOther TaskState = 10
	// StateInternalError should only be given if the task fails prior to execution.
	StateInternalError TaskState = 11
)
