package install

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/cache"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Install install given package
func Install(pkgPath string) error {
	cachePath, err := config.GetCachePath()
	if err != nil {
		return err
	}
	c, err := cache.Read(cachePath)
	if err != nil {
		return err
	}

	pkgName, _, pkgOs, pkgArch, pkgType, err := pkg.ParseFileName(pkgPath)
	if err != nil {
		return err
	}

	// read package
	pkgContent, err := pkg.Read(pkgPath)
	if err != nil {
		return err
	}

	// If binary package override name using alias file
	if pkgType == pkg.Binary {
		if pkgContent["alias"] == nil {
			return fmt.Errorf("no alias file found in package")
		}

		pkgName = string(pkgContent["alias"])
	}

	// Make sure package is not already installed
	if c.GetFiles(pkgName) != nil {
		return fmt.Errorf("package %s is already installed", pkgName)
	}

	files, err := installFromFile(pkgName, pkgOs, pkgArch, pkgType, pkgContent)
	if err != nil {
		return err
	}

	// Everything went well, update local cache
	c.AddPackage(pkgName, files)
	if err := cache.Write(cachePath, c); err != nil {
		return err
	}

	log.Info().Str("package", pkgName).Msg("Successfully installed package")

	return nil
}

func installFromFile(pkgName, pkgOs, pkgArch string, pkgType pkg.Type, pkgContent map[string][]byte) ([]string, error) {
	switch pkgType {
	case pkg.Source:
		files, err := installSourcePackage(pkgContent)
		return files, err
	case pkg.Binary:
		files, err := installBinaryPackage(pkgOs, pkgArch, pkgContent)
		return files, err
	default:
		return nil, fmt.Errorf("can't install package %s", pkgName)
	}
}

func installSourcePackage(pkgContent map[string][]byte) ([]string, error) {
	rootDir, err := config.GetSourceInstallDir()
	if err != nil {
		return nil, err
	}

	var files []string
	for path, content := range pkgContent {
		filePath := filepath.Join(rootDir, path)
		log.Trace().Str("path", filePath).Msg("Writing file")

		// create directory if needed
		if err := os.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
			return nil, err
		}

		// then create file
		if err := ioutil.WriteFile(filePath, content, 0640); err != nil {
			return nil, err
		}

		files = append(files, filePath)
	}

	return files, nil
}

func installBinaryPackage(pkgOs, pkgArch string, pkgContent map[string][]byte) ([]string, error) {
	if pkgOs != runtime.GOOS {
		return nil, fmt.Errorf("package not supported for this os (got: %s want: %s)", pkgOs, runtime.GOOS)
	}

	if pkgArch != runtime.GOARCH {
		return nil, fmt.Errorf("package not supported for this arch (got: %s want: %s)", pkgArch, runtime.GOARCH)
	}

	rootDir, err := config.GetBinaryInstallDir()
	if err != nil {
		return nil, err
	}

	var files []string
	for path, content := range pkgContent {
		if strings.HasPrefix(path, "bin/") {
			realPath := filepath.Join(rootDir, strings.TrimPrefix(path, "bin/"))
			log.Trace().Str("path", realPath).Msg("Writing file")

			// create directory if needed
			if err := os.MkdirAll(filepath.Dir(realPath), 0750); err != nil {
				return nil, err
			}

			if err := ioutil.WriteFile(realPath, content, 0750); err != nil {
				return nil, err
			}

			files = append(files, realPath)
		}
	}

	return files, nil
}
