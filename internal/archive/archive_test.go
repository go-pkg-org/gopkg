package archive

import (
	util "github.com/go-pkg-org/gopkg/internal"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileMapNormal(t *testing.T) {
	cwd, _ := os.Getwd()
	result, _ := CreateFileMap(filepath.Join(cwd, "testfiles"), "some/prefix", []string{})

	expectedPaths := []string{
		filepath.Join(cwd, "testfiles", "test.txt"),
		filepath.Join(cwd, "testfiles", "test.json"),
		filepath.Join(cwd, "testfiles", "test.xml"),
	}
	expectedArchivePaths := []string{
		filepath.Join("some", "prefix", "test.txt"),
		filepath.Join("some", "prefix", "test.json"),
		filepath.Join("some", "prefix", "test.xml"),
	}

	if len(result) != len(expectedPaths) {
		t.Error("length mismatch between expected and result")
	}

	for _, f := range result {
		if !util.Contains(expectedPaths, f.FilePath) {
			t.Errorf("%s did not exist in expected paths", f.FilePath)
		}
		if !util.Contains(expectedArchivePaths, f.ArchivePath) {
			t.Errorf("%s did not exist in expected archive paths", f.ArchivePath)
		}
	}
}

func TestCreateFileMapSpecificType(t *testing.T) {
	cwd, _ := os.Getwd()
	result, _ := CreateFileMap(filepath.Join(cwd, "testfiles"), "some/prefix", []string{".txt", ".json"})

	expectedPaths := []string{
		filepath.Join(cwd, "testfiles", "test.txt"),
		filepath.Join(cwd, "testfiles", "test.json"),
	}
	expectedArchivePaths := []string{
		filepath.Join("some", "prefix", "test.txt"),
		filepath.Join("some", "prefix", "test.json"),
	}

	if len(result) != len(expectedPaths) {
		t.Error("length mismatch between expected and result")
	}

	for _, f := range result {
		if !util.Contains(expectedPaths, f.FilePath) {
			t.Errorf("%s did not exist in expected paths", f.FilePath)
		}
		if !util.Contains(expectedArchivePaths, f.ArchivePath) {
			t.Errorf("%s did not exist in expected archive paths", f.ArchivePath)
		}
	}
}
