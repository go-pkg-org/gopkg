package cmd

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/install"
	"github.com/urfave/cli/v2"
)

// ExecInstall execute the `gopkg install` command
func ExecInstall(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("missing pkg-path")
	}

	return install.Install(c.Args().First())
}
