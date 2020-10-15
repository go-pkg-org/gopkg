package make

import "testing"

func TestGetPackageName(t *testing.T) {
	if getPackageName("github.com/creekorful/mvnparser") != "github-creekorful-mvnparser" {
		t.FailNow()
	}
}

func TestGetVersion(t *testing.T) {
	if v, err := getVersion("vim-2.3.1"); err != nil || v != "2.3.1" {
		t.FailNow()
	}
}
