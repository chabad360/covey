package main

// GetLog reads the unread buffer and adds it to the task's log, then returns that log.
func (t *Task) GetLog() []string {
	// There is a reason why some tty clients print every rewrite of a line,
	// it's much easier to simply process every line as a new line
	// but if we can store the fact that \n hasn't yet been given,
	// we can solve that issue.
	// Also escaping is another issue...
	b := t.Buffer.Bytes()
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
	if len(c) > 0 {
		l = append(l, string(c))
	}
	if len(l) > 0 {
		t.Log = l
	}
	t.Buffer.Reset()
	return t.Log
}
