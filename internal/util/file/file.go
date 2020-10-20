package file

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	errNoFileFound = errors.New("No file found")
)

// FindByExtensions will search for file by provided extensions.
func FindByExtensions(file string, extensions []string) (string, error) {
	path := strings.Replace(file, filepath.Ext(file), "", -1)

	for _, ext := range extensions {
		if ext[0] != '.' {
			ext = "." + ext
		}

		if _, err := os.Stat(path + ext); err == nil {
			return path + ext, nil
		}
	}

	return "", errNoFileFound
}
