package config

import (
	"fmt"
	"os"
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

// GetMaintainerEntry returns the maintainer entry: format Name <Email>
func GetMaintainerEntry() string {
	return fmt.Sprintf("%s <%s>", getEnvOr("GOPKG_MAINTAINER_NAME", "TODO"),
		getEnvOr("GOPKG_MAINTAINER_EMAIL", "TODO"))
}

func getEnvOr(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
