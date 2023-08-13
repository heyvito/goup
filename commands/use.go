package commands

import (
	"fmt"
	"github.com/heyvito/goup/fs"
	"github.com/heyvito/goup/models"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func Use(c models.Context) *cli.Command {
	return &cli.Command{
		Name:        "use",
		Aliases:     []string{"u"},
		Description: "Activates a given version",
		Action: func(context *cli.Context) error {
			if context.NArg() != 1 {
				return cli.Exit("Usage: goup use VERSION", 1)
			}

			installed, err := fs.InstalledVersions(c)
			if err != nil {
				return cli.Exit(err, 1)
			}

			toUse := context.Args().First()

			var versionToUse *models.InstalledVersion
			for _, v := range installed {
				if v.RawName == toUse || v.Version() == toUse {
					versionToUse = &v
					break
				}
			}

			if versionToUse == nil {
				return cli.Exit(fmt.Sprintf("No installed version %s. Use goup list to list installed versions.", toUse), 1)
			}

			currentLink := filepath.Join(c.InstallPath, "current")
			err = os.Remove(currentLink)
			if err != nil && !os.IsNotExist(err) {
				return cli.Exit(err, 1)
			}

			if err = os.Symlink(filepath.Join(versionToUse.Path), currentLink); err != nil {
				fmt.Printf("Failed linking %s to %s:\n", versionToUse.Path, currentLink)
				return cli.Exit(err, 1)
			}

			return nil
		},
	}
}
