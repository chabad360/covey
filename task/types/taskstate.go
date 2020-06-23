package types

const (
	// StateError should be given if the task fails.
	StateError = 0

	// StateDone should be given if the task is complete.
	StateDone = 1

	// StateRunning should be given if the task is running.
	StateRunning = 2

	// StateAborted should be given if the task is aborted (stopped while running).
	StateAborted = 3
	// StateCancelled should be given if the scheduled task is canceled.
	// StateCancelled = 4
	// StateScheduled should be given if the task is scheduled to be executed.
	// StateScheduled = 5

	// StateQueued should be given while the task is waiting to be executed.
	StateQueued = 6

	// StateOther should be avoided.
	StateOther = 10

	// StateInternalError should only be given if the task plugin fails.
	// This will be used if err != nil, if this is given there is a bug in the plugin.
	StateInternalError = 11
)
