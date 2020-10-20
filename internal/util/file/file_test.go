package file

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindByExtensions(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Error(err)
	}

	files := []string{
		filepath.Join(u.HomeDir, ".gopkg.yaml"),
		filepath.Join(u.HomeDir, ".gopkg2.yml"),
		// File without extension can't start with dot.
		filepath.Join(u.HomeDir, "gopkg3"),
	}

	extensions := []string{"yaml", "yml"}

	for _, file := range files {
		path := file

		if filepath.Ext(path) == "" {
			path = path + ".yaml"
		}

		os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)

		out, err := FindByExtensions(file, extensions)
		if err != nil {
			t.Error(err)
		}

		if err := os.Remove(path); err != nil {
			t.Error(err)
		}

		if !strings.Contains(out, file) {
			t.Errorf("File %s not matching output path: %s", file, out)
		}
	}
}
