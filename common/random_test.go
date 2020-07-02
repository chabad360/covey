package common

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	got := RandomString()
	got2 := RandomString()
	if got == got2 {
		t.Errorf("RandomString() returned non random string, %v == %v", got, got2)
	}
}
