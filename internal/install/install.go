package install

import (
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Install install given package
// todo support multiple packages
func Install(pkgPath string) error {
	return installFromFile(pkgPath)
}

// todo add support for named install
func installFromFile(pkgPath string) error {
	log.Info().Str("file", pkgPath).Msg("Installing package from file")
	// todo validate arch + os

	pkgContent, err := pkg.Read(pkgPath)
	if err != nil {
		return err
	}

	// todo centralize logic somewhere to archive package (see: https://github.com/go-pkg-org/gopkg/issues/27)
	if strings.Contains(pkgPath, "-dev") {
		return installSourcePackage(pkgContent)
	}

	// TODO installBinaryPackage
	return nil
}

func installSourcePackage(pkgContent map[string][]byte) error {
	rootDir, err := getSourceInstallDir()
	if err != nil {
		return err
	}

	for path, content := range pkgContent {
		filePath := filepath.Join(rootDir, path)
		log.Debug().Str("path", filePath).Msg("Writing file")

		// create directory if needed
		if err := os.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
			return err
		}

		// then create file
		if err := ioutil.WriteFile(filePath, content, 0640); err != nil {
			return err
		}
	}

	return nil
}

func installBinaryPackage(os, arch string, pkgContent map[string][]byte) error {
	return nil // TODO
}

// getSourceInstallDir returns OS specific installation directory
// for source package
func getSourceInstallDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	// TODO something else? (its enough for now)
	return filepath.Join(u.HomeDir, ".gopkg", "src"), nil
}
