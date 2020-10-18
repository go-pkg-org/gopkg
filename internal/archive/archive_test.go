package archive

import (
	util "github.com/go-pkg-org/gopkg/internal/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateFileMapNormal(t *testing.T) {
	dir, _ := ioutil.TempDir("", "gopkg_*")
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	jsonFile, _ := ioutil.TempFile(dir, "*.json")
	txtFile, _ := ioutil.TempFile(dir, "*.txt")
	xmlFile, _ := ioutil.TempFile(dir, "*.xml")

	result, _ := CreateEntries(dir, "some/prefix", []string{})

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

	result, _ = CreateEntries(dir, "some/prefix", []string{".txt", ".json"})
	expectedPaths = []string{
		filepath.Join(jsonFile.Name()),
		filepath.Join(txtFile.Name()),
	}
	expectedArchivePaths = []string{
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

func TestRead(t *testing.T) {
	dir, _ := ioutil.TempDir("", "gopkg_*")
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	jsonFile, _ := ioutil.TempFile(dir, "*.json")
	jsonFile.WriteString("This is a json file")
	txtFile, _ := ioutil.TempFile(dir, "*.txt")
	txtFile.WriteString("This is a txt file")
	xmlFile, _ := ioutil.TempFile(dir, "*.xml")
	xmlFile.WriteString("This is an xml file")

	xmlFile.Close()
	txtFile.Close()
	jsonFile.Close()

	// Create a tar file to test on in temp dir.
	cmd := exec.Command("tar", "-cf", "out.pkg",
		filepath.Base(jsonFile.Name()),
		filepath.Base(txtFile.Name()),
		filepath.Base(xmlFile.Name()),
	)
	cmd.Dir = dir
	cmd.Run()

	list, err := Read(filepath.Join(dir, "out.pkg"))
	if err != nil {
		t.Errorf("failed to read the archive: %s", err)
	}

	jsonContent := string(list[filepath.Base(jsonFile.Name())])
	xmlContent := string(list[filepath.Base(xmlFile.Name())])
	txtContent := string(list[filepath.Base(txtFile.Name())])

	if strings.EqualFold(jsonContent, "This is a json file") {
		t.Errorf("Json file could not be read.")
	}
	if strings.EqualFold(xmlContent, "This is an xml file") {
		t.Errorf("Xml file could not be read.")
	}
	if strings.EqualFold(txtContent, "This is a txt file") {
		t.Errorf("Txt file could not be read.")
	}

}
