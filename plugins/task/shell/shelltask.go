package main

// GetDetails returns the details of the task.
func (t *Task) GetDetails() map[string]string { return t.Details }

// GetLog reads the unread buffer and adds it to the task's log, then returns that log.
func (t *Task) GetLog() []string {
	if t.Buffer != nil { // Ensure buffer exists
		b := t.Buffer.Bytes()

		var line []byte
		var log []string

		for _, bb := range b { // For each byte...
			if bb == '\n' { // If that byte is a newline:
				log = append(log, string(line)) // Add that line to the log
				line = nil                      // And start the next one
			} else { // Otherwise,
				line = append(line, bb) // Add it to the line
			}
		}

		if len(line) > 0 { // If the last line didn't end with a newline
			log = append(log, string(line)) // Append it
		}

		if len(log) > 0 { // Only set the log if there is stuff on it, otherwise we get empty logs.
			t.Log = log
		}

		t.Buffer.Reset() // Finally, reset the buffer.
	}

	return t.Log
}
