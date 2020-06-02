package main

// GetLog reads the unread buffer and adds it to the task's log, then returns that log.
func (t *Task) GetLog() []string {
	if t.Details.Buffer != nil {
		b := t.Details.Buffer.Bytes()
		c := []byte{}
		l := []string{}
		for _, bb := range b {
			if bb == '\n' {
				l = append(l, string(c))
				c = nil
			} else {
				c = append(c, bb)
			}
		}
		if len(c) > 0 { // Append even if there is no trailing newline.
			l = append(l, string(c))
		}
		if len(l) > 0 { // Set the log to the new value only if there is what to set.
			t.Log = l
		}
		t.Details.Buffer.Reset()
	}
	return t.Log
}
