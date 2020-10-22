package cache

import "testing"

func TestCache(t *testing.T) {
	c := Cache{map[string][]string{}}
	c.AddPackage("gohello", []string{"bin/gohello"})

	if c.GetFiles("gohello")[0] != "bin/gohello" {
		t.Error()
	}

	c.RemovePackage("gohello")

	if c.GetFiles("gohello") != nil {
		t.Error()
	}
}
