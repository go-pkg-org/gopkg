package control

import (
	"fmt"
	"os"
	"path/filepath"
)

const goPkgDir = "gopkg"

func CreateCtrlDirectory(path, version, uploader string, metadata Metadata) error {
	rootDir := filepath.Join(path, goPkgDir)

	if _, err := os.Stat(rootDir); err == nil {
		return fmt.Errorf("%s already exist", rootDir)
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create the root directory
	if err := os.MkdirAll(rootDir, 0750); err != nil {
		return err
	}

	// Create the metadata file
	if err := writeMetadata(metadata, rootDir); err != nil {
		return err
	}

	// Create a default changelog
	if err := writeChangelog(newChangelog(version, uploader), rootDir); err != nil {
		return err
	}

	return nil
}

func ReadCtrlDirectory(path string) (Metadata, Changelog, error) {
	rootDir := filepath.Join(path, goPkgDir)

	m, err := readMetadata(rootDir)
	if err != nil {
		return Metadata{}, Changelog{}, err
	}

	c, err := readChangelog(rootDir)
	if err != nil {
		return Metadata{}, Changelog{}, err
	}

	return m, c, nil
}
