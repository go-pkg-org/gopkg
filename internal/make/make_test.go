package make

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetGitVersion(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "gopkg")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := runGitCmd(tmpDir, nil, "init"); err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("Hello, world!"), 0640); err != nil {
		t.Error(err)
	}

	if err := runGitCmd(tmpDir, nil, "add", "README.md"); err != nil {
		t.Error(err)
	}

	if err := runGitCmd(tmpDir, []string{"GIT_COMMITTER_DATE=\"Thu Oct 15 20:05:34 2020 +0200\""},
		"commit", "-m", "hello"); err != nil {
		t.Error(err)
	}

	v, isTag, err := getGitVersion(tmpDir)
	if err != nil {
		t.Error(err)
	}

	if isTag {
		t.Error("Git version should not be a tag")
	}

	if v != "0.0~git20201015205" {
		t.Errorf("Wrong git version (%s)", v)
	}

	// Create a git tag
	if err := runGitCmd(tmpDir, nil, "tag", "v1.0.0"); err != nil {
		t.Error(err)
	}

	v, isTag, err = getGitVersion(tmpDir)
	if err != nil {
		t.Error(err)
	}

	if !isTag {
		t.Error("Git version should be a tag")
	}

	if v != "v1.0.0" {
		t.Error("Wrong git version")
	}
}

func TestGetMissingDeps(t *testing.T) {
	deps := []string{"github.com/jedib0t/go-pretty/v6/table", "github.com/jedib0t/go-pretty/v6/text",
		"github.com/muesli/termenv", "golang.org/x/crypto/ssh/terminal", "golang.org/x/sys/unix",
		"fmt", "os", "os/exec", "github.com/creekorful/mvnparser/utils"}
	stdDeps := []string{"fmt", "os", "os/exec"}
	importPath := "github.com/creekorful/mvnparser"

	missingDeps, err := getMissingDeps(deps, stdDeps, importPath)
	if err != nil {
		t.Error(err)
	}

	if len(missingDeps) != 4 {
		t.Errorf("Wrong number of missing dependencies found")
	}

	// make sure we've found all dependencies
	depsToFind := map[string]bool{"github.com/muesli/termenv": false,
		"github.com/jedib0t/go-pretty": false, "golang.org/x/crypto": false, "golang.org/x/sys": false}
	for _, d := range missingDeps {
		depsToFind[d] = true
	}
	for d, found := range depsToFind {
		if !found {
			t.Errorf("Missing dep: %s", d)
		}
	}
}

func runGitCmd(dir string, env []string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	if len(env) > 0 {
		cmd.Env = os.Environ()
		for _, val := range env {
			cmd.Env = append(cmd.Env, val)
		}
	}

	b, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error while running `%s` (%s)", cmd.String(), strings.TrimSuffix(string(b), "\n"))
	}

	return nil
}
