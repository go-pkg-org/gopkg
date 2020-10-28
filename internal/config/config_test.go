package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("GOPKG_BIN_DIR", "BIN_DIR")
	os.Setenv("GOPKG_BIN_DIR", "BIN_DIR")

	u, err := user.Current()
	if err != nil {
		t.Error(err)
	}

	path := filepath.Join(u.HomeDir, ".gopkg.yml")
	body := []byte("maintainer:\nemail: test@example.com")

	if err := ioutil.WriteFile(path, body, 0644); err != nil {
		t.Error(err)
	}

	c := &Config{
		Maintainer: Maintainer{
			Name: "Test",
		},
	}

	if err := c.load(); err != nil {
		t.Error(err)
	}

	if c.BinDir != "BIN_DIR" {
		t.Errorf("Config bin dir not equal the expected value, got %s", c.BinDir)
	}

	if c.Maintainer.Name != "Test" {
		t.Errorf("Config maintainer name not equal the expected value, got %s", c.Maintainer.Name)
	}

	if c.Maintainer.Email != "test@example.com" {
		t.Errorf("Config maintainer email not equal the expected value, got %s", c.Maintainer.Email)
	}

	if err := os.Remove(path); err != nil {
		t.Error(err)
	}

	os.Unsetenv("GOPKG_BIN_DIR")
}

func TestConfigMaintainerEnv(t *testing.T) {
	os.Setenv("GOPKG_MAINTAINER_NAME", "Test")

	c := &Config{}

	if err := c.load(); err != nil {
		t.Error(err)
	}

	if c.Maintainer.Name != "Test" {
		t.Errorf("Config maintainer name not equal the expected value, got %s", c.Maintainer.Name)
	}

	os.Unsetenv("GOPKG_MAINTAINER_NAME")
}

func TestConfigDefault(t *testing.T) {
	config, err := Default()
	if err != nil {
		t.Error(err)
	}

	u, err := user.Current()
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		Actual   string
		Expected string
		Text     string
	}{
		{
			Actual:   config.BinDir,
			Expected: filepath.Join(u.HomeDir, GoPkgDir, "bin"),
			Text:     "Default bin dir",
		},
		{
			Actual:   config.CachePath,
			Expected: filepath.Join(u.HomeDir, GoPkgDir, "cache.json"),
			Text:     "Default cache path",
		},
		{
			Actual:   config.SrcDir,
			Expected: filepath.Join(u.HomeDir, GoPkgDir, "src"),
			Text:     "Default src dir",
		},
		{
			Actual:   config.Maintainer.Email,
			Expected: "",
			Text:     "Default maintainer email",
		},
		{
			Actual:   config.Maintainer.Name,
			Expected: "",
			Text:     "Default maintainer name",
		},
	}

	for _, test := range tests {
		if test.Actual != test.Expected {
			t.Errorf("Config %s actual value [%s] is not equal to expected [%s]", test.Text, test.Actual, test.Expected)
		}
	}

	if err := os.Remove(filepath.Join(u.HomeDir, ".gopkg.yaml")); err != nil {
		t.Error(err)
	}
}
