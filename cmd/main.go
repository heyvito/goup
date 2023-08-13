package main

import (
	"github.com/heyvito/goup/commands"
	"github.com/heyvito/goup/models"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	ctx := models.MakeContext()

	app := cli.App{
		Name:      "goup",
		HelpName:  "goup",
		Usage:     "Manages local Golang installations",
		UsageText: "Go-up downloads, installs, and maintains Golang installations based on releases from go.dev",
		Version:   "0.1.0",
		Commands: []*cli.Command{
			commands.List(ctx),
			commands.Info(ctx),
			commands.ListRemote(ctx),
			commands.Install(ctx),
			commands.Use(ctx),
			commands.Uninstall(ctx),
		},
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{Name: "Victor Gama", Email: "hey@vito.io"},
		},
		Copyright: "Copyright (c) 2023 - Victor Gama",
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
