package models

import "testing"

var n = &Node{
	Name: "test",
	ID:   "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1",
}

func TestNode_GetIDShort(t *testing.T) {
	if got := n.GetIDShort(); got != "a7a39b72f29718e6" {
		t.Errorf("Node.GetIDShort() = %v, want %v", got, "a7a39b72f29718e6")
	}
}
