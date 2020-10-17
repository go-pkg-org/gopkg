package cmd

import (
	"fmt"
	make2 "github.com/go-pkg-org/gopkg/internal/make"
	"github.com/urfave/cli/v2"
)

func ExecMake(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("missing import-path")
	}

	return make2.Make(c.Args().First())
}
