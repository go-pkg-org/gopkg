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
)

// Build will build control package located as directory
// and produce binary / dev packages into directory/build folder
func Build(directory string) error {
	m, c, err := control.ReadCtrlDirectory(directory)
	if err != nil {
		return err
	}

	// Recreate build directory
	if err := os.RemoveAll(filepath.Join(directory, "build")); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(directory, "build"), 0750); err != nil {
		return err
	}

	releaseVersion := c.Releases[len(c.Releases)-1].Version

	log.Info().Msgf("Control package: %s %s\n\n", m.Package, releaseVersion)

	// TODO: when supporting dependencies we must fetch it from there
	// and install them

	for _, pkg := range m.Packages {
		var err error
		if pkg.IsSource() {
			if err = buildSourcePackage(directory, releaseVersion, m.ImportPath, pkg); err != nil {
				return err
			}
		} else {
			for targetOs, targetArches := range pkg.Targets {
				for _, targetArch := range targetArches {
					if err = buildBinaryPackage(directory, releaseVersion, targetOs, targetArch, pkg); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func buildSourcePackage(directory, releaseVersion, importPath string, p control.Package) error {
	pkgName, err := pkg.GetName(p.Package, releaseVersion, "", "", true)
	if err != nil {
		return err
	}

	dir, err := pkg.CreateEntries(directory, importPath, []string{})
	if err != nil {
		return err
	}

	if err := pkg.Write(filepath.Join(directory, "build", pkgName), dir, true); err != nil {
		return err
	}

	log.Info().Str("package", p.Package).Msg("Successfully built source package")
	return nil
}

func buildBinaryPackage(directory, releaseVersion, targetOs, targetArch string, p control.Package) error {
	pkgName, err := pkg.GetName(p.Package, releaseVersion, targetOs, targetArch, false)
	if err != nil {
		return err
	}

	buildDir := filepath.Join(directory, "build", pkgName)

	cmd := exec.Command("go", "build", "-o", filepath.Join(buildDir, p.Package), p.Main)
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", targetOs), fmt.Sprintf("GOARCH=%s", targetArch))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Save the package in `build/packageName.pkg`
	err = pkg.Write(filepath.Join(directory, "build", pkgName+".pkg"), []pkg.Entry{
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
