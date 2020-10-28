package main

import (
	"github.com/go-pkg-org/gopkg/internal/pkgarchiver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.DebugLevel)

	app := cli.App{
		Name:    "pkgarchiver",
		Version: "0.0.1",
		Usage:   "Package manager for Golang written applications",
		Authors: []*cli.Author{
			{Name: "Aloïs Micard", Email: "alois@micard.lu"},
			{Name: "Fredrik Forsmo", Email: "hello@frozzare.com"},
			{Name: "Johannes Tegnér", Email: "johannes@jitesoft.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "signing-key",
				Usage: "path to the archive signing key",
			},
			&cli.StringFlag{
				Name:  "maintainer-keyring",
				Usage: "path to the maintainers keyring (to validate incoming package)",
			},
			&cli.StringFlag{
				Name:  "ftp-host",
				Usage: "archive FTP host",
			},
			&cli.StringFlag{
				Name:  "ftp-user",
				Usage: "archive FTP user",
			},
			&cli.StringFlag{
				Name:  "ftp-pass",
				Usage: "archive FTP password",
			},
			&cli.StringFlag{
				Name:  "ftp-dir",
				Usage: "base dir for FTP archive",
			},
		},
		Action: pkgarchiver.Execute,
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Msg("error while running application")
		os.Exit(1)
	}
}
