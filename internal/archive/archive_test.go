package archive

import (
	util "github.com/go-pkg-org/gopkg/internal/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileMapNormal(t *testing.T) {
	dir, _ := ioutil.TempDir("", "*")
	defer os.RemoveAll(dir)

	jsonFile, _ := ioutil.TempFile(dir, "*.json")
	txtFile, _ := ioutil.TempFile(dir, "*.txt")
	xmlFile, _ := ioutil.TempFile(dir, "*.xml")

	result, _ := CreateFileMap(dir, "some/prefix", []string{})

	expectedPaths := []string{
		filepath.Join(jsonFile.Name()),
		filepath.Join(txtFile.Name()),
		filepath.Join(xmlFile.Name()),
	}
	expectedArchivePaths := []string{
		filepath.Join("some", "prefix", filepath.Base(txtFile.Name())),
		filepath.Join("some", "prefix", filepath.Base(jsonFile.Name())),
		filepath.Join("some", "prefix", filepath.Base(xmlFile.Name())),
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
	dir, _ := ioutil.TempDir("", "*")
	defer os.RemoveAll(dir)

	jsonFile, _ := ioutil.TempFile(dir, "*.json")
	txtFile, _ := ioutil.TempFile(dir, "*.txt")
	ioutil.TempFile(dir, "*.xml")

	result, _ := CreateFileMap(dir, "some/prefix", []string{".txt", ".json"})

	expectedPaths := []string{
		filepath.Join(jsonFile.Name()),
		filepath.Join(txtFile.Name()),
	}
	expectedArchivePaths := []string{
		filepath.Join("some", "prefix", filepath.Base(txtFile.Name())),
		filepath.Join("some", "prefix", filepath.Base(jsonFile.Name())),
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
