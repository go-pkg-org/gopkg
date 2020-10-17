package build

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/control"
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

	fmt.Printf("Control package: %s %s\n\n", m.Package, releaseVersion)

	// TODO: when supporting dependencies we must fetch it from there
	// and install them

	for _, pkg := range m.Packages {
		fmt.Printf("Building %s\n", pkg.Package)

		var err error
		if pkg.IsSource() {
			err = buildSourcePackage(directory, releaseVersion, pkg)
		} else {
			err = buildBinaryPackage(directory, releaseVersion, pkg)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func buildSourcePackage(directory, releaseVersion string, pkg control.Package) error {
	pkgName := fmt.Sprintf("%s_%s.pkg", pkg.Package, releaseVersion)
	cmd := exec.Command("tar", "-czvf", "build/"+pkgName, ".")
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func buildBinaryPackage(directory, releaseVersion string, pkg control.Package) error {
	// TODO support multi arch / os
	pkgName := fmt.Sprintf("%s_%s_darwin_amd64", pkg.Package, releaseVersion)
	buildDir := fmt.Sprintf("build/%s", pkgName)

	cmd := exec.Command("go", "build", "-v", "-o", fmt.Sprintf("%s/usr/share/bin/%s", buildDir, pkg.Package), pkg.Main)
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// finally create tar archive
	cmd = exec.Command("tar", "-czvf", "build/"+pkgName+".pkg", "-C", buildDir, ".") // TODO: strip '.'
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Successfully build %s.pkg\n", pkgName)

	return nil
}
