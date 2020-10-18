package main

import (
	"github.com/go-pkg-org/gopkg/internal/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.DebugLevel)

	app := cli.App{
		Name:  "gopkg",
		Usage: "Package manager for Golang written applications",
		Authors: []*cli.Author{
			{"Aloïs Micard", "alois@micard.lu"},
			{"Fredrik Forsmo", "hello@frozzare.com"},
			{"Johannes Tegnér", "johannes@jitesoft.com"},
		},
		Commands: []*cli.Command{
			{
				Name:      "make",
				Usage:     "create a new package from import-path",
				ArgsUsage: "import-path",
				Action:    cmd.ExecMake,
			},
			{
				Name:      "build",
				Usage:     "build a package from control directory",
				ArgsUsage: "control-directory",
				Action:    cmd.ExecBuild,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Msg("error while running application")
		os.Exit(1)
	}
}
