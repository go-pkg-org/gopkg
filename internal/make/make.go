package make

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/go-pkg-org/gopkg/internal/control"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/go-pkg-org/gopkg/internal/util"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Make create a brand new control package from given import path
func Make(importPath string) error {
	pkgName := pkg.GetName(importPath)

	if _, err := os.Stat(pkgName); err == nil {
		return fmt.Errorf("already existing package directory: %s", pkgName)
	}

	// Fetch & extract upstream source code
	version, err := getUpstreamSource(importPath, pkgName)
	if err != nil {
		return err
	}
	// Remove any leading v since we doesn't want it in gopkg archive
	cleanVersion := strings.TrimPrefix(version, "v")

	// Get defined importPaths (dependencies)
	deps, err := getImportPaths(pkgName)
	if err != nil {
		return err
	}

	// Get std dependencies (builtin)
	stdDeps, err := getStdDeps()
	if err != nil {
		return err
	}

	// Then get its dependencies
	missingDeps, err := getMissingDeps(deps, stdDeps, importPath)
	if err != nil {
		return err
	}

	if len(missingDeps) > 0 {
		log.Warn().Strs("dependencies", missingDeps).Msg("Dependencies that need to be packaged first")
		return nil
	}

	m := control.Metadata{
		Maintainers: []string{config.GetMaintainerEntry()},
		Packages:    []control.Package{},
		ImportPath:  importPath,
	}

	// Search for binary packages
	binPkgs, err := getBinaryPackages(importPath, pkgName)
	if err != nil {
		return err
	}
	m.Packages = append(m.Packages, binPkgs...)

	// Create the control directory
	if err := control.CreateCtrlDirectory(pkgName, cleanVersion, config.GetMaintainerEntry(), m); err != nil {
		return err
	}

	log.Info().
		Str("import-path", importPath).
		Str("version", cleanVersion).
		Str("package", pkgName).
		Msg("Detected values")
	for _, p := range m.Packages {
		log.Info().Str("package", p.Alias).Msg("Built package")
	}

	return nil
}

// Get the package missing dependencies dependencies
// - remove the 'std' dependencies (builtin)
// - remove the dependencies that belongs to the project we want to package
// todo remove already packaged deps
func getMissingDeps(deps, stdDeps []string, importPath string) ([]string, error) {
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

		// Filter dep to make sure we're only adding 'root' import path for the dependencies
		// and also prevent duplicates
		rootDep := ""
		parts := strings.Split(dep, "/")
		if len(parts) == 3 {
			rootDep = dep
		} else {
			rootDep = fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
		}

		// And make sure it's not already packaged
		if !util.Contains(realDeps, rootDep) {
			realDeps = append(realDeps, rootDep)
		}
	}

	return realDeps, nil
}

// Get the package define import paths
func getImportPaths(path string) ([]string, error) {
	cmd := exec.Command("go", "list", "-f", "'{{ join .Imports \"\\n\" }}'", "./...")
	cmd.Dir = path
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

// getExecutables will lookup for executable in given directory and returns their corresponding package
func getBinaryPackages(importPath, directory string) ([]control.Package, error) {
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

		if strings.Contains(string(b), "func main()") && strings.HasSuffix(path, ".go") {
			fileName := strings.Replace(info.Name(), ".go", "", 1)
			aliasName := filepath.Join(importPath, fileName)
			pkgs = append(pkgs, control.Package{
				Alias:       aliasName,
				Description: "TODO",
				Main:        strings.TrimPrefix(path, directory+"/"),
				BinName:     strings.ReplaceAll(aliasName, "/", "-"),
				Targets:     getDefaultTargets(),
			})
			log.Trace().Str("file", path).Str("alias", aliasName).Msg("Found binary package")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return pkgs, nil
}

// getUpstreamSource fetch latest available upstream source
// this method return path to upstream source, version, and error if any
func getUpstreamSource(importPath, where string) (string, error) {
	remote := fmt.Sprintf("https://%s.git", importPath)
	log.Debug().Str("remote", remote).Msg("Found upstream remote")

	// Clone repository
	log.Debug().Str("remote", remote).Msg("Cloning remote")
	cmd := exec.Command("git", "clone", remote, where)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error: git clone %s result (%s)", remote, err)
	}

	// Get git repository latest version
	version, isTag, err := getGitVersion(where)
	if err != nil {
		return "", err
	}
	log.Debug().Str("version", version).Bool("tagged", isTag).Msg("Found upstream version")

	// if this is a tagged release, checkout it to align source code
	if isTag {
		log.Debug().Str("tag", version).Msg("Checking out tag")
		cmd = exec.Command("git", "checkout", version)
		cmd.Dir = where
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}

	return version, nil
}

// getGitVersion will attempt to auto-detect the latest stable/tagged release
// if upstream tag release: it will return the latest tag
// if upstream doesn't tag release: it will create a special version for the latest (HEAD) commit
func getGitVersion(gitDir string) (string, bool, error) {
	// Extract latest tag / version
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = gitDir
	b, err := cmd.Output()
	if err != nil {
		// There maybe no tag available, create a manual version using commit date
		cmd = exec.Command("git", "--no-pager", "log", "-1", "--date=short", "--pretty=format:%cD")
		cmd.Dir = gitDir

		b, err = cmd.Output()
		if err != nil {
			// were doomed
			return "", false, err
		}

		date, err := time.Parse(time.RFC1123Z, strings.TrimSuffix(string(b), "\n"))
		if err != nil {
			return "", false, err
		}

		return fmt.Sprintf("0.0~git%d%d%d%d%d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute()), false, nil
	}

	return strings.TrimSuffix(string(b), "\n"), true, nil
}

func getDefaultTargets() map[string][]string {
	return map[string][]string{
		"linux":  {"amd64"},
		"darwin": {"amd64"},
	}
}
