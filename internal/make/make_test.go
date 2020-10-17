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

func TestGetPackageName(t *testing.T) {
	if getPackageName("github.com/creekorful/mvnparser") != "github-creekorful-mvnparser" {
		t.FailNow()
	}
}

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