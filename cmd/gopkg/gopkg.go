package main

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/build"
	make2 "github.com/go-pkg-org/gopkg/internal/make"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		_, _ = fmt.Fprintf(os.Stderr, "correct usage: gopkg <action> args...")
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
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
