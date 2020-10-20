package build

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/control"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/rs/zerolog/log"
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

	m, c, err := control.ReadCtrlDirectory(path)
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

	releaseVersion := c.Releases[len(c.Releases)-1].Version

	log.Info().Msgf("Control package: %s %s", m.Package, releaseVersion)

	// TODO: when supporting dependencies we must fetch it from there
	// and install them

	for _, p := range m.Packages {
		var err error
		if p.IsSource() {
			if err = buildSourcePackage(path, releaseVersion, m.ImportPath, p); err != nil {
				return err
			}
		} else {
			for targetOs, targetArches := range p.Targets {
				for _, targetArch := range targetArches {
					if err = buildBinaryPackage(path, releaseVersion, targetOs, targetArch, p); err != nil {
						return err
					}
				}
			}
		}
	}

	// Finally build control package
	return buildControlPackage(path, m.Package, releaseVersion)
}

func extractControlPackage(path string) (string, error) {
	// Make sure its a control package
	fileName := filepath.Base(path)
	_, _, _, _, pkgType, err := pkg.ParseName(fileName)
	if err != nil {
		return "", err
	}
	if pkgType != pkg.Control {
		return "", fmt.Errorf("%s is not a control package", fileName)
	}

	content, err := pkg.Read(path)
	if err != nil {
		return "", err
	}

	baseDir := filepath.Dir(path)
	for p, b := range content {
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

func buildControlPackage(directory, pkgName string, releaseVersion string) error {
	pkgName, err := pkg.GetName(pkgName, releaseVersion, "", "", pkg.Control)
	if err != nil {
		return err
	}

	dir, err := pkg.CreateEntries(directory, strings.TrimSuffix(pkgName, "."+pkg.FileExt), []string{".git"})
	if err != nil {
		return err
	}

	// Save the package in `./<pkgName>`
	if err := pkg.Write(pkgName, dir, true); err != nil {
		return err
	}

	log.Info().Str("package", pkgName).Msg("Successfully built control package")
	return nil
}

func buildSourcePackage(directory, releaseVersion, importPath string, p control.Package) error {
	pkgName, err := pkg.GetName(p.Package, releaseVersion, "", "", pkg.Source)
	if err != nil {
		return err
	}

	dir, err := pkg.CreateEntries(directory, importPath, []string{".git", control.GoPkgDir})
	if err != nil {
		return err
	}

	// Save the package in `./<pkgName>`
	if err := pkg.Write(pkgName, dir, true); err != nil {
		return err
	}

	log.Info().Str("package", pkgName).Msg("Successfully built source package")
	return nil
}

func buildBinaryPackage(directory, releaseVersion, targetOs, targetArch string, p control.Package) error {
	pkgName, err := pkg.GetName(p.Package, releaseVersion, targetOs, targetArch, pkg.Binary)
	if err != nil {
		return err
	}

	buildDir := filepath.Join(directory, "build", pkgName)

	cmd := exec.Command("go", "build", "-o", filepath.Join(buildDir, p.Package), p.Main)
	log.Trace().Msgf("Executing `%s`", cmd.String())
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", targetOs), fmt.Sprintf("GOARCH=%s", targetArch))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Save the package in `./<pkgName>`
	err = pkg.Write(filepath.Join(pkgName), []pkg.Entry{
		{
			FilePath:    filepath.Join(buildDir, p.Package),
			ArchivePath: filepath.Join("bin", p.Package),
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
