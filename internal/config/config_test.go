package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("GOPKG_BIN_DIR", "BIN_DIR")

	u, err := user.Current()
	if err != nil {
		t.Error(err)
	}

	file := filepath.Join(u.HomeDir, ".gopkg.yml")
	body := []byte(`maintainer:
  email: test@example.com`)

	if err := ioutil.WriteFile(file, body, 0644); err != nil {
		t.Error(err)
	}

	c := &Config{
		Maintainer: Maintainer{
			Name: "Test",
		},
	}

	if err := c.Load(); err != nil {
		fmt.Println(err)
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

	if err := os.Remove(file); err != nil {
		t.Error(err)
	}
}
