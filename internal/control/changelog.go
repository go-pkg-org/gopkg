package control

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const changelogFile = "changelog.yaml"

// Changelog is the root object containing all package releases
type Changelog struct {
	Releases []Release
}

// Release is produced each time a package is released
type Release struct {
	// The package version number (upstream-internal)
	// f.e 1.2.0-1 is the initial release of upstream version 1.2.0.
	Version string
	// Who has taking care of the release upload
	Uploader string
	// The human descriptions of changes applied since last release
	Changes []string
}

// NewChangelog create a brand new changelog
func newChangelog(initialVersion, uploader string) Changelog {
	return Changelog{
		Releases: []Release{{
			Version:  fmt.Sprintf("%s-1", initialVersion),
			Uploader: uploader,
			Changes:  []string{"Initial release"},
		}},
	}
}

// WriteChangelog write the given changelog
func writeChangelog(c Changelog, path string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(path, changelogFile), b, 0640)
}

// ReadChangelog read changelog from file
func readChangelog(path string) (Changelog, error) {
	var c Changelog

	f, err := os.Open(filepath.Join(path, changelogFile))
	if err != nil {
		return Changelog{}, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return Changelog{}, err
	}

	return c, nil
}
