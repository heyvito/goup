package commands

import (
	"fmt"
	"github.com/heyvito/goup/fs"
	"github.com/heyvito/goup/godev"
	"github.com/heyvito/goup/models"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
	"strings"
)

func ClearLine() {
	fmt.Printf("\r\u001B[2K")
}

func ListRemote(c models.Context) *cli.Command {
	return &cli.Command{
		Name:    "list-remote",
		Aliases: []string{"lr"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Obtains all versions, including archived and unstable",
			},
		},
		Description: "Lists remote installable SDKs from go.dev",
		Action: func(context *cli.Context) error {
			fmt.Printf("Querying go.dev...")
			list, err := godev.ListVersions(context.Bool("all"))

			if err != nil {
				ClearLine()
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			installed, err := fs.InstalledVersions(c)
			ClearLine()

			for _, i := range list {
				ver := i.PrettyVersion()
				fmt.Printf(" - %s", ver)
				if i.Stable {
					fmt.Printf(" (stable, ")
				} else {
					fmt.Printf(" (unstable, ")
				}
				for _, v := range installed {
					if v.Version() == ver {
						fmt.Printf("installed, ")
						break
					}
				}
				compatible := false
				archs := i.Archs()
				currentArch := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
				for _, v := range archs {
					if v == currentArch {
						compatible = true
						break
					}
				}
				if compatible {
					fmt.Println("compatible)")
				} else {
					fmt.Println("incompatible)")
				}
				fmt.Printf("    Archs: %s\n", strings.Join(archs, ", "))
			}
			return nil
		},
	}
}
