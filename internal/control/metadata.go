package control

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const metadataFile = "metadata.yaml"

// Metadata represent the package metadata
type Metadata struct {
	// Package is the control package name
	Package string
	// List of the package maintainers
	// i.e who take the responsibility for uploading & managing it
	Maintainers []string
	// List of the packages built by this control package
	Packages []Package
}

// Package represent a package installable
type Package struct {
	// The package name
	Package string
	// List of architectures for which the package should be built
	Architectures []string
}

// WriteMetadata write the given metadata
func WriteMetadata(m Metadata, path string) error {
	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", path, metadataFile), b, 0640)
}
