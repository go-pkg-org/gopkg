package build

import (
	"fmt"
	util "github.com/go-pkg-org/gopkg/internal"
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
	pkgName := fmt.Sprintf("%s_%s.pkg", pkg.Package, releaseVersion)

/*	cmd := exec.Command("tar", "-czvf", "build/"+pkgName, ".")
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard*/

	dir, err := util.CreateFileMap(directory)
	if err != nil {
		return err
	}

	if err := util.CreateTar(filepath.Join("build", pkgName), dir, true); err != nil {
		return err
	}

	/*if err := cmd.Run(); err != nil {
		return err
	}*/

	return nil
}

func buildBinaryPackage(directory, releaseVersion, targetOs, targetArch string, pkg control.Package) error {
	pkgName := fmt.Sprintf("%s_%s_%s_%s", pkg.Package, releaseVersion, targetOs, targetArch)
	buildDir := fmt.Sprintf("build/%s", pkgName)

	// TODO: /usr/share/bin is specific to linux/darwin arch. We should have specialized path depending on targetOs
	cmd := exec.Command("go", "build", "-v", "-o", fmt.Sprintf("%s/usr/share/bin/%s", buildDir, pkg.Package), pkg.Main)
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", targetOs), fmt.Sprintf("GOARCH=%s", targetArch))
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
