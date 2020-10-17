package main

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/build"
	make2 "github.com/go-pkg-org/gopkg/internal/make"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
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
		err = build.Build(os.Args[2])
	default:
		err = fmt.Errorf("unknow action")
	}

	if err != nil {
		log.Err(err).Msg("error while running gopkg")
		os.Exit(1)
	}
}
