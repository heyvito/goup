package commands

import (
	"fmt"
	"github.com/heyvito/goup/fs"
	"github.com/heyvito/goup/models"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func Uninstall(c models.Context) *cli.Command {
	return &cli.Command{
		Name:        "uninstall",
		Usage:       "uninstall VERSION",
		Description: "Removes a given Go version",
		Action: func(context *cli.Context) error {
			version := context.Args().First()
			if version == "" {
				fmt.Printf("Usage: goup uninstall VERSION. Use goup list to obtain a list of installed versions.\n")
				os.Exit(1)
			}

			installed, err := fs.InstalledVersions(c)
			if err != nil {
				return cli.Exit(err, 1)
			}

			var toRemove *models.InstalledVersion
			for _, v := range installed {
				if v.RawName == version || v.Version() == version {
					toRemove = &v
					break
				}
			}

			active, err := fs.CurrentVersion(c)
			if err != nil {
				return cli.Exit(err, 1)
			}

			if toRemove == nil {
				return cli.Exit(fmt.Errorf("unknown version %s. Use goup list to obtain a list of installed versions", version), 1)
			}

			if toRemove.Ok {
				fmt.Printf("You are about to uninstall %s, which is marked as healthy. Would you like to continue? (y/N) ", toRemove.RawName)
				r, err := ReadChar()
				if err != nil {
					return cli.Exit(err, 1)
				}
				if strings.ToLower(r) != "y" {
					fmt.Println("Aborted.")
					os.Exit(1)
				}
			}

			fmt.Printf("Removing %s...", toRemove.RawName)
			if err = os.RemoveAll(toRemove.Path); err != nil {
				fmt.Println(" Failed.")
				return cli.Exit(err, 1)
			}
			fmt.Println(" OK.")
			if active != nil && toRemove.Path == active.Path {
				fmt.Printf("WARNING: You just removed the active go version. Use goup use VERSION to change to another version.\n")
			}
			return nil
		},
	}
}
