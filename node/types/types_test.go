package types

import "testing"

var n = &Node{
	Name:    "test",
	ID:      "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1",
	Plugin:  "plugin",
	Details: "test",
}

func TestNode_GetName(t *testing.T) {
	if got := n.GetName(); got != "test" {
		t.Errorf("Node.GetName() = %v, want %v", got, "test")
	}
}

func TestNode_GetID(t *testing.T) {
	if got := n.GetID(); got != "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1" {
		t.Errorf("Node.GetID() = %v, want %v", got, "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1")
	}
}

func TestNode_GetIDShort(t *testing.T) {
	if got := n.GetIDShort(); got != "a7a39b72f29718e6" {
		t.Errorf("Node.GetIDShort() = %v, want %v", got, "a7a39b72f29718e6")
	}
}

func TestNode_GetPlugin(t *testing.T) {
	if got := n.GetPlugin(); got != "plugin" {
		t.Errorf("Node.GetPlugin() = %v, want %v", got, "plugin")
	}
}

func TestNode_GetDetails(t *testing.T) {
	if got := n.GetDetails(); got != "test" {
		t.Errorf("Node.GetDetails() = %v, want %v", got, "test")
	}
}
