package install

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// Install install given package
// todo support multiple packages
func Install(pkgPath string) error {
	pkgName, err := installFromFile(pkgPath)
	if err != nil {
		return err
	}
	log.Info().Str("package", pkgName).Msg("Successfully installed package")

	return nil
}

// todo add support for named install
func installFromFile(pkgPath string) (string, error) {
	pkgName := filepath.Base(pkgPath)
	pkgName, _, pkgOs, pkgArch, isSrc, err := pkg.ParseName(pkgName)
	if err != nil {
		return pkgName, err
	}

	log.Info().Str("package", pkgName).Msg("Installing package")

	pkgContent, err := pkg.Read(pkgPath)
	if err != nil {
		return pkgName, err
	}

	if isSrc {
		return pkgName, installSourcePackage(pkgContent)
	}

	return pkgName, installBinaryPackage(pkgOs, pkgArch, pkgContent)
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

func installBinaryPackage(pkgOs, pkgArch string, pkgContent map[string][]byte) error {
	if pkgOs != runtime.GOOS {
		return fmt.Errorf("package not supported for this os (got: %s want: %s)", pkgOs, runtime.GOOS)
	}

	if pkgArch != runtime.GOARCH {
		return fmt.Errorf("package not supported for this arch (got: %s want: %s)", pkgArch, runtime.GOARCH)
	}

	rootDir, err := getBinaryInstallDir()
	if err != nil {
		return err
	}

	for path, content := range pkgContent {
		if strings.HasPrefix(path, "bin/") {
			realPath := filepath.Join(rootDir, strings.TrimPrefix(path, "bin/"))
			log.Debug().Str("path", realPath).Msg("Writing file")

			// create directory if needed
			if err := os.MkdirAll(filepath.Dir(realPath), 0750); err != nil {
				return err
			}

			if err := ioutil.WriteFile(realPath, content, 0750); err != nil {
				return err
			}
		}
	}

	return nil
}

// getBinaryInstallDir returns OS specific installation directory for bin package
func getBinaryInstallDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	// TODO something else? (its enough for now)
	return filepath.Join(u.HomeDir, ".gopkg", "bin"), nil
}

// getSourceInstallDir returns OS specific installation directory for source package
func getSourceInstallDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	// TODO something else? (its enough for now)
	return filepath.Join(u.HomeDir, ".gopkg", "src"), nil
}
