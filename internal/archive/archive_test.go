package archive

import (
	util "github.com/go-pkg-org/gopkg/internal/util"
	"io/ioutil"
	"os"
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
	jsonFile.Close()
	txtFile, _ := ioutil.TempFile(dir, "*.txt")
	txtFile.WriteString("This is a txt file")
	txtFile.Close()
	xmlFile, _ := ioutil.TempFile(dir, "*.xml")
	xmlFile.WriteString("This is an xml file")
	xmlFile.Close()

	err := Write(filepath.Join(dir, "out.pkg"), []Entry{
		{xmlFile.Name(), "test/xmlfile.xml"},
		{jsonFile.Name(), "jsonfile.json"},
		{txtFile.Name(), "txtfile.txt"},
	}, true)

	if err != nil {
		t.Errorf("failed to create archive: %s", err)
	}

	list, err := Read(filepath.Join(dir, "out.pkg"))
	if err != nil {
		t.Errorf("failed to read the archive: %s", err)
	}

	jsonContent := string(list["jsonfile.json"])
	xmlContent := string(list["test/xmlfile.xml"])
	txtContent := string(list["txtfile.txt"])

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
