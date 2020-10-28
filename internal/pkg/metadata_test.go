package pkg

import "testing"

func TestMeta_Source(t *testing.T) {
	m := Meta{
		Alias:       "trandoshan/crawler",
		Main:        "crawler.go",
		BinName:     "tdsh-crawler",
		Description: "",
		Targets:     map[string][]string{},
	}

	if m.IsSource() {
		t.FailNow()
	}

	m = Meta{
		Alias:       "github.com/creekorful/mvnparser",
		Main:        "",
		BinName:     "",
		Description: "",
		Targets:     map[string][]string{},
	}

	if !m.IsSource() {
		t.FailNow()
	}
}
