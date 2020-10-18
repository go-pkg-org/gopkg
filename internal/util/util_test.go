package util

import "testing"

func TestContains(t *testing.T) {
	slice := []string{
		"abc", "efg", "111", "222", "333",
	}

	if !Contains(slice, "efg") {
		t.Errorf("failed to find string in slice")
	}

	if Contains(slice, "hhh") {
		t.Errorf("contains returned true with none existing value")
	}
}
