package config

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/go-pkg-org/gopkg/internal/util/file"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const configFile = ".gopkg.yaml"

// Maintainer is the object containg info about the maintainer.
type Maintainer struct {
	Email string `yaml:"email"`
	Name  string `yaml:"name"`
}

// Config is the root object containg the configuration file.
type Config struct {
	BinDir     string     `yaml:"bin_dir" envconfig:"bin_dir"`
	Maintainer Maintainer `yaml:"maintainer"`
	SrcDir     string     `yaml:"src_dir"  envconfig:"src_dir"`
}

// Load loads the configuration file from the users home directory.
func (c *Config) Load() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	path, err := file.FindByExtensions(filepath.Join(u.HomeDir, configFile), []string{"yaml", "yml"})
	if err == nil {
		out, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal([]byte(out), &c)
		if err != nil {
			return err
		}
	}

	err = envconfig.Process("gopkg", c)
	if err != nil {
		return err
	}

	return nil
}
