package make

import "testing"

func TestGetPackageName(t *testing.T) {
	if getPackageName("github.com/creekorful/mvnparser") != "github-creekorful-mvnparser" {
		t.FailNow()
	}
}
