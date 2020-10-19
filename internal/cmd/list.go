package cmd

import (
	"github.com/go-pkg-org/gopkg/internal/list"
	"github.com/urfave/cli/v2"
)

func ExecList(c *cli.Context) error {
	return list.List(c.Bool("installed"))
}
