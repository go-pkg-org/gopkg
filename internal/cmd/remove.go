package cmd

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/remove"
	"github.com/urfave/cli/v2"
)

func ExecRemove(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("missing pkg-name")
	}

	return remove.Remove(c.Args().First())
}
