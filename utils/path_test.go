package utils

import "testing"

func TestRuntimeDir(t *testing.T) {
	t.Log(RuntimeDir(".ansible"))
}
