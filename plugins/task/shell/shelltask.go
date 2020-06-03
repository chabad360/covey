package main

// GetDetails returns the details of the task.
func (t *Task) GetDetails() interface{} { return t.Details }

// GetLog reads the unread buffer and adds it to the task's log, then returns that log.
func (t *Task) GetLog() []string {
	if t.Details.Buffer != nil {
		b := t.Details.Buffer.Bytes()
		line := []byte(nil)
		log := []string(nil)
		for _, bb := range b {
			if bb == '\n' {
				log = append(log, string(line))
				line = nil
			} else {
				line = append(line, bb)
			}
		}
		if len(line) > 0 { // Append even if there is no trailing newline.
			log = append(log, string(line))
		}
		if len(log) > 0 { // Set the log to the new value only if there is what to set.
			t.Log = log
		}
		t.Details.Buffer.Reset()
	}
	return t.Log
}
