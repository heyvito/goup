package commands

import (
	"fmt"
	"github.com/heyvito/goup/fs"
	"github.com/heyvito/goup/models"
	"github.com/urfave/cli/v2"
	"os"
)

func List(c models.Context) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Aliases:     []string{"l"},
		Description: "Lists local installed SDKs managed by go-up",
		Action: func(context *cli.Context) error {
			exists, err := fs.DirExists(c.InstallPath)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			if !exists {
				fmt.Printf("There are no SDK installed through go-up.\n")
				fmt.Printf("To install a new one, use goup install latest")
				return nil
			}

			list, err := fs.InstalledVersions(c)
			if err != nil {
				fmt.Printf("Failed querying installed versions: %s\n", err)
				return nil
			}

			if len(list) == 0 {
				fmt.Printf("No versions installed locally. Use list-remote to list installable versions.\n")
			}

			active, err := fs.CurrentVersion(c)
			if err != nil {
				fmt.Printf("Failed querying current installed version: %s\n", err)
				return nil
			}

			for _, i := range list {
				fmt.Printf(" - %s", i.Version())
				if active != nil && i.Path == active.Path {
					fmt.Printf(" active")
				}
				if i.Ok {
					fmt.Printf(" healthy")
				} else {
					fmt.Printf(" degraded: %s", i.Status)
				}
				fmt.Println()
			}

			return nil
		},
	}
}
