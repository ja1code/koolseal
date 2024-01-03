package main

import (
	"log"
	"os"

	"github.com/ja1code/koolseal/commands"
	"github.com/urfave/cli/v2"
)

var (
	version string
)

func main() {
	app := &cli.App{
		Name:     "Koolseal",
		Version:  version,
		Usage:    "A easier way to manage Kubeseal secrets",
		Commands: commands.Commands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
