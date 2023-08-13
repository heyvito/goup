package commands

import (
	"bytes"
	"fmt"
	"github.com/heyvito/goup/fs"
	"github.com/heyvito/goup/godev"
	"github.com/heyvito/goup/models"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ReadChar() (string, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return "", err
	}
	fmt.Printf(string(b[0]) + "\r\n")
	return string(b[0]), nil
}

func Install(c models.Context) *cli.Command {
	return &cli.Command{
		Name:        "install",
		Usage:       "install VERSION",
		Description: "Installs a given Go version",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "non-interactive",
				Aliases: []string{"n"},
				Usage:   "Do not ask for confirmation when installing",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Automatically overwrite existing installations, even in non-interactive mode.",
			},
		},
		Action: func(context *cli.Context) error {
			version := context.Args().First()
			if version == "" {
				fmt.Printf("Usage: goup install VERSION. Use goup list-remote to obtain a list of availabel versions.\n")
				fmt.Printf("Alternatively, you can also pass 'latest' as the version.\n")
				os.Exit(1)
			}

			fmt.Printf("Querying go.dev...")
			versions, err := godev.ListVersions(true)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			installed, err := fs.InstalledVersions(c)
			if err != nil {
				fmt.Println(" Failed.")
				fmt.Printf("Error checking installed versions: %s\n", err)
				os.Exit(1)
			}

			if strings.ToLower(version) == "latest" {
				version = versions[0].Version
			}

			var toInstall *models.RemoteVersion
			for _, v := range versions {
				if v.Version == version || v.PrettyVersion() == version {
					toInstall = &v
					break
				}
			}

			fmt.Println(" OK.")

			if toInstall == nil {
				fmt.Printf("Error: Version %s not found. Use goup list-remote to obtain a list of available versions.\n", version)
				os.Exit(1)
			}

			nonInteractive := context.Bool("non-interactive")
			force := context.Bool("force")

			// Is it compatible with the current system?
			compatible := false
			currentPlatform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
			for _, v := range toInstall.Archs() {
				if v == currentPlatform {
					compatible = true
					break
				}
			}

			if !nonInteractive {
				fmt.Printf("About to download and install %s. Continue? (y/N) ", toInstall.Version)
				option, err := ReadChar()
				if err != nil {
					fmt.Println()
					fmt.Printf("Error reading from standard input: %s\n", err)
					os.Exit(1)
				}
				if strings.ToLower(option) != "y" {
					fmt.Printf("Aborted.")
					os.Exit(1)
				}
			}

			if !compatible {
				fmt.Printf("%s is not compatible with your platform. Aborting.", toInstall.Version)
				os.Exit(1)
			}

			for _, v := range installed {
				if v.RawName != toInstall.Version {
					continue
				}

				if !v.Ok {
					fmt.Printf("WARNING: Overwriting degraded installed version %s\n", toInstall.Version)
				} else if nonInteractive && !force {
					fmt.Printf("Refusing to overwrite healthy version %s. Use -f to force.\n", toInstall.Version)
					os.Exit(1)
				} else if force {
					fmt.Printf("WARNING: Overwriting healthy version %s.\n", toInstall.Version)
				} else if !nonInteractive {
					fmt.Printf("It seems there is a healthy version of %s already installed. Would you like to reinstall it? (y/N) ", toInstall.Version)
					option, err := ReadChar()
					if err != nil {
						fmt.Println()
						fmt.Printf("Error reading from standard input: %s\n", err)
						os.Exit(1)
					}
					if strings.ToLower(option) != "y" {
						fmt.Printf("Aborted.")
						os.Exit(1)
					}
				}
			}

			fmt.Println()

			var file *models.RemoteVersionFile
			for _, v := range toInstall.Files {
				if v.OS == runtime.GOOS &&
					v.Arch == runtime.GOARCH &&
					v.Kind == "archive" {
					file = &v
					break
				}
			}

			if file == nil {
				fmt.Println("ERROR: Could not find target version on files list for the selected version.")
				fmt.Println("This is likely a bug. Please report it.")
				os.Exit(1)
			}

			fmt.Printf("Downloading version %s...", toInstall.Version)
			path, err := godev.DownloadVersion("https://go.dev/dl/" + file.Filename)
			if err != nil {
				fmt.Printf(" Failed.\n%s\n", err)
				return cli.Exit("Failed downloading archive", 1)
			}

			defer func() { _ = os.Remove(path) }()

			fmt.Printf(" OK.\n")
			fmt.Printf("Checking file integrity... ")
			ok, err := godev.CheckShasum(path, file.Sha256)
			if err != nil {
				fmt.Println(" Failed.")
				return cli.Exit(err, 1)
			}

			if !ok {
				fmt.Println(" Failed.")
				fmt.Println("Signatures did not match.")
				return cli.Exit("Download might be corrupted. Please try again.", 1)
			}

			fmt.Println("OK.")
			fmt.Printf("Preparing to decompress archive...")
			targetPath := filepath.Join(c.InstallPath, "versions", toInstall.Version)
			if err = os.RemoveAll(targetPath); err != nil {
				fmt.Printf(" Failed.")
				return cli.Exit(err, 1)
			}

			if err = os.MkdirAll(targetPath, 0755); err != nil {
				fmt.Println(" Failed.")
				return cli.Exit(err, 1)
			}

			fmt.Println("OK.")
			fmt.Printf("Decompressing SDK...")

			if err = decompress(path, targetPath); err != nil {
				fmt.Println(" Failed.")
				return cli.Exit(err, 1)
			}

			fmt.Println(" OK.")

			fmt.Printf("Checking installation...")
			cmd := exec.Command(filepath.Join(targetPath, "bin", "go"), "version")
			outBuf := &bytes.Buffer{}
			errBuf := &bytes.Buffer{}
			cmd.Stdout = outBuf
			cmd.Stderr = errBuf
			if err = cmd.Start(); err != nil {
				fmt.Println(" Failed.")
				if rmErr := os.RemoveAll(targetPath); rmErr != nil {
					fmt.Printf("Also: Failed removing broken installation at %s: %s", targetPath, rmErr)
				}
				return cli.Exit(err, 1)
			}

			if err = cmd.Wait(); err != nil {
				fmt.Println(" Failed.")
				fmt.Printf("%s\n%s\n", outBuf.String(), errBuf.String())
				if rmErr := os.RemoveAll(targetPath); rmErr != nil {
					fmt.Printf("Also: Failed removing broken installation at %s: %s", targetPath, rmErr)
				}
				return cli.Exit(err, 1)
			}

			fmt.Printf(" OK. %s", outBuf.String())
			fmt.Printf("Instalation completed. Use goup use %s to activate the new version.\n", toInstall.PrettyVersion())

			checkEnv(c)
			return nil
		},
	}
}

func checkEnv(c models.Context) {
	pathPresent := false
	dirPath := filepath.Join(c.InstallPath, "bin")
	for _, v := range strings.Split(os.Getenv("PATH"), ":") {
		if v == dirPath {
			pathPresent = true
			break
		}
	}

	pathOk := os.Getenv("GOPATH") == c.InstallPath

	if !pathPresent {
		fmt.Printf("WARNING: %s is not present in your path. Make sure to include it before other directories.\n", dirPath)
	}

	if !pathOk {
		fmt.Printf("WARNING: It is advised to point your GOPATH to %s as well. Update your environments :)\n", c.InstallPath)
	}

}

func decompress(tarGzPath, dest string) error {
	cmd := exec.Command("tar", "--strip-components=1", "-xvf", tarGzPath)
	cmd.Dir = dest
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("could not decompress archive: %s\n%s\nError: %w", outBuf.String(), errBuf.String(), err)
	}

	return nil
}

func copyClose(dst io.WriteCloser, src io.Reader) error {
	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		return err
	}
	return dst.Close()
}
