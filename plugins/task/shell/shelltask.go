package main

// GetLog reads the unread buffer and adds it to the task's log, then returns that log.
func (t *Task) GetLog() []string {
	// As of now whenever we get a CR not followed by LF, we reset the line, just like a tty
	// if there is a reason to change this, it can be.
	if t.Details.Buffer != nil {
		b := t.Details.Buffer.Bytes()
		c := []byte{}
		l := []string{}
		cr := false
		for _, bb := range b {
			if bb == '\n' {
				l = append(l, string(c))
				c = nil
				cr = false
			} else if bb == '\r' {
				c = append(c, bb)
				cr = true
			} else {
				if cr == true {
					c = nil
				}
				c = append(c, bb)
			}
		}
		if len(c) > 0 {
			l = append(l, string(c))
		}
		if len(l) > 0 {
			t.Log = l
		}
		t.Details.Buffer.Reset()
	}
	return t.Log
}
