package cmd

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/build"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func ExecBuild(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("missing control-directory")
	}

	realPath := c.Args().First()
	if !filepath.IsAbs(realPath) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		realPath = filepath.Join(wd, realPath)
	}

	return build.Build(realPath)
}
