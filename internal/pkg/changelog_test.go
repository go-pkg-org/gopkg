package pkg

import "testing"

func TestChangelog_LastRelease(t *testing.T) {
	c := Changelog{Releases: []Release{}}
	c.Releases = append(c.Releases, Release{Version: "1.0.0-1"}, Release{Version: "1.0.0-2"})

	lastRelease, err := c.LastRelease()
	if err != nil {
		t.FailNow()
	}

	if lastRelease.Version != "1.0.0-2" {
		t.Errorf("wrong last release")
	}
}

func TestNewChangelog(t *testing.T) {
	c := newChangelog("1.0.0", "Aloïs Micard <alois@micard.lu>")
	if len(c.Releases) != 1 {
		t.FailNow()
	}

	if got := c.Releases[0].Version; got != "1.0.0-1" {
		t.Errorf("wrong release (got %s)", got)
	}
	if got := c.Releases[0].Uploader; got != "Aloïs Micard <alois@micard.lu>" {
		t.Errorf("wrong uploader (got %s)", got)
	}
	if got := len(c.Releases[0].Changes); got != 1 {
		t.Errorf("wrong number of changes (got %d)", got)
	}
	if got := c.Releases[0].Changes[0]; got != "Initial packaging" {
		t.Errorf("wrong changes (got %s)", got)
	}
}
