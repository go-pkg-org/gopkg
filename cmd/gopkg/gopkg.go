package main

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/build"
	make2 "github.com/go-pkg-org/gopkg/internal/make"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.DebugLevel)

	if len(os.Args) == 1 {
		log.Error().Msg("correct usage: gopkg <action> args...")
		os.Exit(1)
	}



	var err error
	action := os.Args[1]
	switch action {
	case "make":
		err = make2.Make(os.Args[2])
	case "build":
		realPath := os.Args[2]
		if !filepath.IsAbs(realPath) {
			wd, err := os.Getwd()
			if err != nil {
				log.Panic().Msg("failed to get working directory")
			}
			realPath = filepath.Join(wd, realPath)
		}

		err = build.Build(realPath)
	default:
		err = fmt.Errorf("unknow action")
	}

	if err != nil {
		log.Err(err).Msg("error while running gopkg")
		os.Exit(1)
	}
}
