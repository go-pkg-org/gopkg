package build

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/archive"
	"github.com/go-pkg-org/gopkg/internal/control"
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
		log.Info().Str("package", pkg.Package).Msg("Building source package")

		var err error
		if pkg.IsSource() {
			if err = buildSourcePackage(directory, releaseVersion, pkg); err != nil {
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

func buildSourcePackage(directory, releaseVersion string, pkg control.Package) error {
	pkgName := fmt.Sprintf("%s_%s-dev.pkg", pkg.Package, releaseVersion)

	dir, err := archive.CreateFileMap(directory, "", []string{})
	if err != nil {
		return err
	}

	if err := archive.Create(filepath.Join("build", pkgName), dir, true); err != nil {
		return err
	}

	log.Info().Str("package", pkg.Package).Msg("Successfully built package")
	return nil
}

func buildBinaryPackage(directory, releaseVersion, targetOs, targetArch string, pkg control.Package) error {
	pkgName := fmt.Sprintf("%s_%s_%s_%s", pkg.Package, releaseVersion, targetOs, targetArch)
	buildDir := filepath.Join("build", pkgName)

	cmd := exec.Command("go", "build", "-v", "-o", filepath.Join(buildDir, pkg.Package), pkg.Main)
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", targetOs), fmt.Sprintf("GOARCH=%s", targetArch))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Save the package in `build/packageName.pkg`
	err := archive.Create(filepath.Join("build", pkgName+".pkg"), []archive.Entry{
		{
			FilePath:    filepath.Join(buildDir, pkg.Package),
			ArchivePath: filepath.Join("bin", pkg.Package),
		},
	}, true)

	if err != nil {
		return err
	}

	// Remove the build file and keep package.
	if err := os.RemoveAll(filepath.Join(buildDir)); err != nil {
		return err
	}

	log.Info().Str("package", pkg.Package).Msg("Successfully built package")
	return nil
}
