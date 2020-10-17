package build

import (
	"archive/tar"
	"bytes"
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



	cmd := exec.Command("tar", "-czvf", "build/"+pkgName, ".")
	cmd.Dir = directory
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return err
	}

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

// Create a tar file from a set of files.
// The expected input `files` should be 'filename  => filePath' where the filename
// is the path that will be created inside the tar file, and the filePath is the path
// to the file to add.
// The path is the path to the where the tar file will be created.
func createTar(path string, files map[string]string, overwrite bool) error {
	if !overwrite {
		if _, err := ioutil.ReadFile(path); err != nil {
			return fmt.Errorf("failed to create new tar source. File already exist")
		}
	}

	var buffer bytes.Buffer
	tw := tar.NewWriter(&buffer)

	for fileName, file := range files {
		fileBody, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		header := &tar.Header{
			Name: fileName,
			Mode: 0644,
			Size: int64(len(fileBody)),
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := tw.Write(fileBody); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return ioutil.WriteFile(path, buffer.Bytes(), 0755)
}