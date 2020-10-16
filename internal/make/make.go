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
		Maintainers: []string{getMaintainerEntry()},
		Packages: []control.Package{
			// Create initial source package
			{Package: pkgName + "-dev", Description: "Package development files"},
		},
	}

	// Search for binary packages
	binPkgs, err := getBinaryPackages(directory)
	if err != nil {
		return err
	}
	m.Packages = append(m.Packages, binPkgs...)

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

// getExecutables will lookup for executable in given directory and returns their corresponding package
func getBinaryPackages(directory string) ([]control.Package, error) {
	var pkgs []control.Package

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

		// todo better lookup
		if strings.Contains(string(b), "func main()") && strings.HasSuffix(path, ".go") {
			pkgs = append(pkgs, control.Package{
				Package:       strings.Replace(info.Name(), ".go", "", 1),
				Description:   "TODO",
				Main:          strings.TrimPrefix(path, directory+"/"),
				Architectures: []string{"amd64"},
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return pkgs, nil
}

// getMaintainerEntry returns the maintainer entry: format Name <Email>
func getMaintainerEntry() string {
	return fmt.Sprintf("%s <%s>", getEnvOr("GOPKG_MAINTAINER_NAME", "TODO"),
		getEnvOr("GOPKG_MAINTAINER_EMAIL", "TODO"))
}

func getEnvOr(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
