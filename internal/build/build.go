package build

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Build will build control package located as directory
// and produce binary / dev packages into directory/build folder
func Build(path string) error {
	// If path is pointing to a .pkg file, extract it
	if strings.HasSuffix(path, "."+pkg.FileExt) {
		log.Debug().Str("package", path).Msg("Extracting control package")

		p, err := extractControlPackage(path)
		if err != nil {
			return err
		}
		path = p
	}

	m, c, err := pkg.ReadCtrlDirectory(path)
	if err != nil {
		return err
	}

	// Recreate build directory
	if err := os.RemoveAll(filepath.Join(path, "build")); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(path, "build"), 0750); err != nil {
		return err
	}

	goPath, err := config.GetGoPathDir()
	if err != nil {
		return err
	}

	// Get latest release
	latestRelease, err := c.LastRelease()
	if err != nil {
		return err
	}
	releaseVersion := latestRelease.Version

	log.Info().
		Str("importPath", m.ImportPath).
		Str("version", releaseVersion).
		Msgf("Building for control package")

	// Run unit tests
	cmd := exec.Command("go", "test", "./...")
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%s", goPath))
	cmd.Stderr = os.Stderr
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return err
	}

	// Build source package
	if err := buildSourcePackage(path, m.ImportPath, releaseVersion); err != nil {
		return err
	}

	for _, p := range m.Packages {
		for targetOs, targetArches := range p.Targets {
			for _, targetArch := range targetArches {
				if err = buildBinaryPackage(goPath, path, releaseVersion, targetOs, targetArch, p); err != nil {
					return err
				}
			}
		}
	}

	// Finally build control package
	return buildControlPackage(path, m.ImportPath, releaseVersion)
}

func extractControlPackage(path string) (string, error) {
	// Make sure its a control package
	fileName := filepath.Base(path)
	// TODO refactor this last usage of ParseFileName
	_, _, _, _, pkgType, err := pkg.ParseFileName(fileName)
	if err != nil {
		return "", err
	}
	if pkgType != pkg.Control {
		return "", fmt.Errorf("%s is not a control package", fileName)
	}

	f, err := pkg.ReadFile(path)
	if err != nil {
		return "", err
	}

	baseDir := filepath.Dir(path)
	for p, b := range f.Files() {
		targetPath := filepath.Join(baseDir, p)
		log.Debug().Str("path", targetPath).Msg("Writing file")

		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(targetPath), 0750); err != nil {
			return "", err
		}

		// Then creating file
		if err := ioutil.WriteFile(targetPath, b, 0640); err != nil {
			return "", err
		}
	}

	return strings.TrimSuffix(path, "."+pkg.FileExt), nil
}

func buildControlPackage(directory, importPath string, releaseVersion string) error {
	fileName, err := pkg.GetFileName(importPath, releaseVersion, "", "", pkg.Control)
	if err != nil {
		return err
	}

	dir, err := pkg.CreateEntries(directory, strings.TrimSuffix(fileName, "."+pkg.FileExt), []string{".git"})
	if err != nil {
		return err
	}

	// Save the package in `./<fileName>`
	if err := pkg.Write(fileName, dir, true); err != nil {
		return err
	}

	log.Info().Str("package", fileName).Msg("Successfully built control package")
	return nil
}

func buildSourcePackage(directory, importPath, releaseVersion string) error {
	fileName, err := pkg.GetFileName(importPath, releaseVersion, "", "", pkg.Source)
	if err != nil {
		return err
	}

	dir, err := pkg.CreateEntries(directory, importPath, []string{".git", pkg.GoPkgDir})
	if err != nil {
		return err
	}

	// Create package definition
	p := pkg.Meta{
		Alias:          importPath,
		ReleaseVersion: releaseVersion,
	}
	b, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	// TODO improve that
	if err := ioutil.WriteFile(filepath.Join(directory, "package.yaml"), b, 0640); err != nil {
		return err
	}
	dir = append(dir, pkg.Entry{
		FilePath:    filepath.Join(directory, "package.yaml"),
		ArchivePath: "package.yaml",
	})

	// Save the package in `./<fileName>`
	if err := pkg.Write(fileName, dir, true); err != nil {
		return err
	}

	log.Info().Str("package", fileName).Msg("Successfully built source package")
	return nil
}

func buildBinaryPackage(goPath, directory, releaseVersion, targetOs, targetArch string, p pkg.Meta) error {
	pkgName, err := pkg.GetFileName(p.Alias, releaseVersion, targetOs, targetArch, pkg.Binary)
	if err != nil {
		return err
	}

	buildDir := filepath.Join(directory, "build", pkgName)

	cmd := exec.Command("go", "build", "-o", filepath.Join(buildDir, p.BinName), p.Main)
	log.Trace().Msgf("Executing `%s`", cmd.String())
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", targetOs),
		fmt.Sprintf("GOARCH=%s", targetArch), fmt.Sprintf("GOPATH=%s", goPath))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Copy package definition
	// this is used later on to determinate which package we are installing & populate info
	p.Targets = nil
	p.TargetOS = targetOs
	p.TargetArch = targetArch
	p.ReleaseVersion = releaseVersion
	b, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(buildDir, "package.yaml"), b, 0640); err != nil {
		return err
	}

	// Save the package in `./<pkgName>`
	err = pkg.Write(filepath.Join(pkgName), []pkg.Entry{
		// Add the binary
		{
			FilePath:    filepath.Join(buildDir, p.BinName),
			ArchivePath: filepath.Join("bin", p.BinName),
		},
		// Add the alias file
		{
			FilePath:    filepath.Join(buildDir, "package.yaml"),
			ArchivePath: "package.yaml",
		},
	}, true)

	if err != nil {
		return err
	}

	// Remove the build file and keep package.
	if err := os.RemoveAll(filepath.Join(buildDir)); err != nil {
		return err
	}

	log.Info().Str("package", pkgName).Msg("Successfully built binary package")
	return nil
}
