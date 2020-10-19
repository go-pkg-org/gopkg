package config

import (
	"os/user"
	"path/filepath"
)

// GetBinaryInstallDir returns OS specific installation directory for bin package
func GetBinaryInstallDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	// TODO something else? (its enough for now)
	return filepath.Join(u.HomeDir, ".gopkg", "bin"), nil
}

// GetSourceInstallDir returns OS specific installation directory for source package
func GetSourceInstallDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	// TODO something else? (its enough for now)
	return filepath.Join(u.HomeDir, ".gopkg", "src"), nil
}

// GetCachePath returns path to the installed package cache
func GetCachePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, ".gopkg", "cache.json"), nil
}
