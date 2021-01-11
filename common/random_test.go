package common_test

import (
	"testing"

	"github.com/chabad360/covey/common"
)

func TestRandomString(t *testing.T) {
	got := common.RandomString()
	got2 := common.RandomString()
	if got == got2 {
		t.Errorf("RandomString() not random, %v == %v", got, got2)
	}
}
