package make

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/control"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Make(directory string) error {
	// Acquire temp go path for the session
	goPath, err := getTempGoPath()
	if err != nil {
		return err
	}

	// Get import path
	importPath, err := getImportPath(directory)
	if err != nil {
		return err
	}

	// Get version
	version, err := getVersion(directory)
	if err != nil {
		return err
	}

	// Then get its dependencies
	deps, err := getRealDeps(goPath, directory)
	if err != nil {
		return err
	}

	pkgName := getPackageName(importPath)

	if len(deps) > 0 {
		fmt.Printf("Dependencies that need to be packaged first:\n")
		for _, dep := range deps {
			fmt.Printf("- %s\n", dep)
		}
		return nil // TODO error instead?
	}

	m := control.Metadata{
		Package:     pkgName,
		Maintainers: []string{"Aloïs Micard <alois@micard.lu>"},
		Packages: []control.Package{
			// Create initial source package
			{Package: pkgName + "-dev"},
		},
	}

	// Search for executables
	executables, err := getExecutables(directory)
	if err != nil {
		return err
	}
	// and create package for them
	for _, executable := range executables {
		m.Packages = append(m.Packages, control.Package{
			Package:       executable,
			Architectures: []string{"all"},
		})
	}

	// Create the control directory
	if err := control.CreateCtrlDirectory(directory, version, "Aloïs Micard <alois@micard.lu>", m); err != nil {
		return err
	}

	fmt.Printf("Import-Path: %s\n", importPath)
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Control package: %s\n", pkgName)
	fmt.Printf("Built packages:\n")
	for _, pkg := range m.Packages {
		fmt.Printf("-> %s\n", pkg.Package)
	}

	return nil
}

// Get the package 'real' dependencies
// - remove the 'std' dependencies (builtin)
// - remove the dependencies that belongs to the project we want to package
func getRealDeps(goPath, importPath string) ([]string, error) {
	// Then get its dependencies
	deps, err := getDeps(goPath, importPath)
	if err != nil {
		return nil, err
	}

	// Get std dependencies
	stdDeps, err := getStdDeps()
	if err != nil {
		return nil, err
	}

	// Get 'real' dependencies
	// i.e exclude std dependencies
	var realDeps []string
	for _, dep := range deps {
		// Ignore internal dependencies
		if strings.HasPrefix(dep, importPath) {
			continue
		}

		// Ignore std dependencies
		isStdDep := false
		for _, stdDep := range stdDeps {
			if dep == stdDep {
				isStdDep = true
				break
			}
		}
		if isStdDep {
			continue
		}

		realDeps = append(realDeps, dep)
	}

	return realDeps, nil
}

// Get the package define dependencies
func getDeps(goPath, path string) ([]string, error) {
	cmd := exec.Command("go", "list", "-f", "'{{ join .Imports \"\\n\" }}'", "./...")
	cmd.Dir = path
	cmd.Env = getEnvVariables(goPath)
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseLines(b), nil
}

// Get the builtin dependencies
func getStdDeps() ([]string, error) {
	cmd := exec.Command("go", "list", "std")
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseLines(b), nil
}

func parseLines(b []byte) []string {
	// sanitize output
	output := strings.ReplaceAll(string(b), "'", "")

	// remove trailing \n
	output = strings.TrimSuffix(output, "\n")

	return strings.Split(output, "\n")
}

// Get env variables to use when creating go subprocess
func getEnvVariables(goPath string) []string {
	return []string{
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("GOPATH=%s", goPath),
		"GO111MODULE=off",
		"PATH=/usr/local/bin:/usr/bin:/bin",
	}
}

func getImportPath(directory string) (string, error) {
	cmd := exec.Command("go", "list")
	cmd.Dir = directory
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(b), "\n"), nil
}

func getTempGoPath() (string, error) {
	return "/tmp/gopkg-gopath", nil
}

// Retrieve the package version using the directory name
func getVersion(directory string) (string, error) {
	if strings.Count(directory, "-") != 1 {
		return "", fmt.Errorf("malformed directory name %s. expected <name>-<version>", directory)
	}

	return strings.Split(directory, "-")[1], nil
}

// Translate from importPath to package name
// i.e github.com/creekorful/mvnparser -> github-creekorful-mvnparser
func getPackageName(importPath string) string {
	return strings.Replace(strings.ReplaceAll(importPath, "/", "-"), "github.com", "github", 1)
}

// getExecutables will lookup for executable in given directory and returns their name
func getExecutables(directory string) ([]string, error) {
	var executables []string

	if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(string(b), "func main()") {
			executables = append(executables, strings.Replace(info.Name(), ".go", "", 1))
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return executables, nil
}
