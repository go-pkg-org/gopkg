package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-pkg-org/gopkg/internal/util/file"
	"gopkg.in/yaml.v2"
)

const metadataFile = "metadata.yaml"

// ControlMeta represent the control package metadata
type ControlMeta struct {
	// The Go import path
	ImportPath string
	// List of the package maintainers
	// i.e who take the responsibility for uploading & managing it
	Maintainers []string
	// The package build dependencies (i.e what we need to pull before building the package)
	BuildDependencies []string `yaml:"build_dependencies"`
	// List of the packages built by this control package
	Packages []Meta
}

// Meta represent a package information
type Meta struct {
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
	// These fields below are copied into the package.yaml definition
	TargetOS       string `yaml:"target_os,omitempty"`
	TargetArch     string `yaml:"target_arch,omitempty"`
	ReleaseVersion string `yaml:"release_version,omitempty"`
}

// IsSource determinate if package is a source one
func (m *Meta) IsSource() bool {
	return m.Main == ""
}

func writeControlMeta(m ControlMeta, path string) error {
	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", path, metadataFile), b, 0640)
}

func readControlMetadata(path string) (ControlMeta, error) {
	var m ControlMeta

	path, err := file.FindByExtensions(filepath.Join(path, metadataFile), []string{"yaml", "yml"})
	if err != nil {
		return ControlMeta{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		return ControlMeta{}, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&m); err != nil {
		return ControlMeta{}, err
	}

	return m, nil
}
