package main

import (
	"log"
	"os"

	"github.com/ja1code/koolseal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: commands.Commands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
