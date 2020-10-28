package file

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	_, _ = rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

func TestFindByExtensions(t *testing.T) {
	files := []string{
		tempFileName(".", ".yaml"),
		tempFileName(".", ".yml"),
		// File without extension can't start with dot.
		tempFileName("", ""),
	}

	extensions := []string{"yaml", "yml"}

	for _, file := range files {
		path := file

		if filepath.Ext(path) == "" {
			path = path + ".yaml"
		}

		_, _ = os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)

		out, err := FindByExtensions(file, extensions)
		if err != nil {
			t.Error(err)
		}

		if runtime.GOOS != "windows" {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}

		if !strings.Contains(out, file) {
			t.Errorf("File %s not matching output path: %s", file, out)
		}
	}
}
