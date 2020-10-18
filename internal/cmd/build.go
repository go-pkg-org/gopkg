package cmd

import (
	"github.com/go-pkg-org/gopkg/internal/build"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

// ExecBuild execute the `gopkg build` command
func ExecBuild(c *cli.Context) error {
	path := c.Args().First()
	if path == "" {
		path = "."
	}

	absolutePath, err := getAbsolutePath(path)
	if err != nil {
		return err
	}

	return build.Build(absolutePath)
}

func getAbsolutePath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(wd, path), nil
	}

	return path, nil
}
