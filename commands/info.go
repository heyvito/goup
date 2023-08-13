package commands

import (
	"fmt"
	"github.com/heyvito/goup/fs"
	"github.com/heyvito/goup/models"
	"github.com/urfave/cli/v2"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Info(c models.Context) *cli.Command {
	return &cli.Command{
		Name:        "info",
		Aliases:     []string{"i"},
		Description: "Shows current system configuration",
		Action: func(context *cli.Context) error {
			goos := runtime.GOOS
			arch := runtime.GOARCH
			goBinPath, goBinPathErr := exec.LookPath("go")

			fmt.Printf("OS: %s\n", goos)
			fmt.Printf("Architecture: %s\n", arch)
			if goBinPathErr != nil {
				fmt.Printf("Current Go binary path: Not found (%s)\n", goBinPathErr)
			} else {
				goBinPath, goBinPathErr = filepath.EvalSymlinks(goBinPath)
				if goBinPathErr != nil {
					fmt.Printf("Current Go binary path: Not found (%s)\n", goBinPathErr)
				} else {
					fmt.Printf("Current Go binary path: %s\n", goBinPath)
				}
			}

			fmt.Printf("Go-up data path: %s\n", c.InstallPath)
			fmt.Printf("Go-up installed SDKs:\n")
			versions, err := fs.InstalledVersions(c)
			if err != nil {
				fmt.Printf("  Error checking: %s\n", err)
			} else {
				if len(versions) == 0 {
					fmt.Printf("  No SDKs installed. Install using goup install.\n")
				}
				for _, v := range versions {
					fmt.Printf("  - %s: ", v.Version())
					if v.Ok {
						fmt.Printf("OK\n")
						continue
					}

					fmt.Printf("Possibly invalid:\n")
					fmt.Printf("    %s\n", v.Status)
				}
			}

			current, err := fs.CurrentVersion(c)
			fmt.Printf("Active installation:\n")
			if err != nil {
				fmt.Printf("  Error: %s\n", err)
			} else {
				if current == nil {
					fmt.Printf("  No installation is currently active.\n")
				} else {
					fmt.Printf("  Version: %s\n", current.Version())
					fmt.Printf("  Health: ")
					if current.Ok {
						fmt.Printf("Ok.\n")
					} else {
						fmt.Printf("Degraded:\n")
						fmt.Printf("    %s\n", current.Status)
					}
				}
			}

			if goBinPathErr == nil {
				fmt.Printf("System is using go-up binary? ")
				if current == nil {
					fmt.Printf("No.\n")
				} else if strings.HasPrefix(goBinPath, current.Path) {
					fmt.Printf("Yes. (%s)\n", goBinPath)
				} else {
					fmt.Printf("No. (%s)\n", goBinPath)
				}
			}

			return nil
		},
	}
}
