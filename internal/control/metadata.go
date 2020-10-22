package control

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const metadataFile = "metadata.yaml"

// Metadata represent the package metadata
type Metadata struct {
	// The Go import path
	ImportPath string
	// List of the package maintainers
	// i.e who take the responsibility for uploading & managing it
	Maintainers []string
	// The package build dependencies (i.e what we need to pull before building the package)
	BuildDependencies []string `yaml:"build_dependencies"`
	// List of the packages built by this control package
	Packages []Package
}

// Package represent a package installable
type Package struct {
	// The package alias (i.e what the user will use to identify the package)
	Alias string
	// Main is the full path to the file containing the `func main()` sequence
	Main string `yaml:"main,omitempty"`
	// BinName is the name of the binary that will be installed
	BinName string
	// Human description of the package
	Description string
	// Targets describe the build target (os,arches)
	Targets map[string][]string `yaml:"targets,omitempty"`
}

// writeMetadata write the given metadata
func writeMetadata(m Metadata, path string) error {
	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", path, metadataFile), b, 0640)
}

// ReadMetadata read metadata from file
func readMetadata(path string) (Metadata, error) {
	var m Metadata

	f, err := os.Open(filepath.Join(path, metadataFile))
	if err != nil {
		return Metadata{}, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&m); err != nil {
		return Metadata{}, err
	}

	return m, nil
}
