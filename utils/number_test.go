package utils

import "testing"

func TestRandomNumber(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Log(RandomNumber("", 8))
	}
}
