package make

import (
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

	if err := runGitCmd(tmpDir, "init"); err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("Hello, world!"), 0640); err != nil {
		t.Error(err)
	}

	if err := runGitCmd(tmpDir, "add", "README.md"); err != nil {
		t.Error(err)
	}

	if err := runGitCmd(tmpDir, "commit", "-m", "hello"); err != nil {
		t.Error(err)
	}

	v, isTag, err := getGitVersion(tmpDir)
	if err != nil {
		t.Error(err)
	}

	if isTag {
		t.Error("Git version should not be a tag")
	}

	if !strings.HasPrefix(v, "0.0~git") {
		t.Error("Wrong git version")
	}

	// Create a git tag
	if err := runGitCmd(tmpDir, "tag", "v1.0.0"); err != nil {
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

func runGitCmd(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}
