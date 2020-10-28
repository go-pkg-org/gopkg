package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/archive"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrPackageAlreadyInstalled is returns when the package we are trying to install is already installed
var ErrPackageAlreadyInstalled = errors.New("package is already installed")

// ErrWrongTarget is returned when the package we are trying to install is not compatible
var ErrWrongTarget = errors.New("package is not compatible")

// Cache is a local gopkg cache
type Cache interface {
	InstallPkgFile(filePath string) (pkg.Meta, error)
	InstallPkg(aliasName string) (pkg.Meta, error)
	ListPackages(onlyInstalled bool) ([]string, error)
	RemovePkg(alias string) error
}

type cache struct {
	Packages  map[string][]string `json:"packages"`
	arcClient archive.Client
	cacheFile string
	conf      *config.Config
}

func (c *cache) InstallPkgFile(filePath string) (pkg.Meta, error) {
	p, err := pkg.ReadFile(filePath)
	if err != nil {
		return pkg.Meta{}, nil
	}

	// Try to install package
	meta, err := c.installPkg(p)
	if err != nil {
		return pkg.Meta{}, fmt.Errorf("error while installing package %s: %s", filePath, err)
	}

	return meta, nil
}

func (c *cache) InstallPkg(aliasName string) (pkg.Meta, error) {
	p, err := c.arcClient.GetLatestRelease(aliasName, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return pkg.Meta{}, nil
	}

	return c.installPkg(p)
}

func (c *cache) installPkg(pkgFile pkg.File) (pkg.Meta, error) {
	// Read meta file
	meta, err := pkgFile.Metadata()
	if err != nil {
		return pkg.Meta{}, err
	}

	// Make sure package is not already installed
	if _, exist := c.Packages[meta.Alias]; exist {
		return pkg.Meta{}, ErrPackageAlreadyInstalled
	}

	var files []string
	// source package can be installed no matter what
	if meta.IsSource() {
		files, err = installSourcePkg(pkgFile, c.conf.SrcDir)
		if err != nil {
			return pkg.Meta{}, err
		}
	} else {
		// binary package need to match os / arch
		if meta.TargetOS != runtime.GOOS || meta.TargetArch != runtime.GOARCH {
			return pkg.Meta{}, ErrWrongTarget
		}

		files, err = installBinaryPkg(pkgFile, c.conf.BinDir)
		if err != nil {
			return pkg.Meta{}, err
		}
	}

	// Update local cache
	c.Packages[meta.Alias] = files
	if err := write(c.cacheFile, c); err != nil {
		return pkg.Meta{}, err
	}

	return meta, err
}

func installSourcePkg(pkgFile pkg.File, sourceInstallDir string) ([]string, error) {
	var files []string
	for path, content := range pkgFile.Files() {
		// Do not install package.yaml or package.yml file
		if path == "package.yaml" || path == "package.yml" {
			continue
		}

		filePath := filepath.Join(sourceInstallDir, path)
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

func installBinaryPkg(pkgFile pkg.File, binaryInstallDir string) ([]string, error) {
	var files []string
	for path, content := range pkgFile.Files() {
		// Do not install package.yaml or package.yml file
		if path == "package.yaml" || path == "package.yml" {
			continue
		}

		if strings.HasPrefix(path, "bin/") {
			realPath := filepath.Join(binaryInstallDir, strings.TrimPrefix(path, "bin/"))
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

func (c *cache) ListPackages(onlyInstalled bool) ([]string, error) {
	var pkgs []string
	if onlyInstalled {
		for p := range c.Packages {
			pkgs = append(pkgs, p)
		}
	} else {
		idx, err := c.arcClient.GetIndex()
		if err != nil {
			return nil, err
		}

		for p := range idx.Packages {
			pkgs = append(pkgs, p)
		}
	}

	return pkgs, nil
}

func (c *cache) RemovePkg(alias string) error {
	files, exist := c.Packages[alias]
	if !exist {
		return fmt.Errorf("package %s not installed", alias)
	}

	// remove installed files
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			log.Warn().Str("file", file).Str("err", err.Error()).Msg("unable to delete file")
		}
	}

	// update cache
	delete(c.Packages, alias)

	return write(c.cacheFile, c)
}

// NewCache create a brand new cache using given arguments
func NewCache(cacheFile string, arcClient archive.Client, conf *config.Config) (Cache, error) {
	c, err := read(cacheFile)
	if err != nil {
		return nil, err
	}

	c.arcClient = arcClient
	c.cacheFile = cacheFile
	c.conf = conf

	return c, nil
}

func read(path string) (*cache, error) {
	var c cache
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &cache{Packages: map[string][]string{}}, nil
		}
		return nil, err
	}

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Write a cache to target path
func write(path string, cache Cache) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(cache); err != nil {
		return err
	}
	return nil
}
