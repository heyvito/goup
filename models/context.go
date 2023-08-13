package models

import (
	"fmt"
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
)

type Context struct {
	InstallPath string
}

func MakeContext() Context {
	ctx := Context{}
	if v, ok := os.LookupEnv("GOUP_INSTALL_DIR"); ok {
		ctx.InstallPath = v
	} else {
		ctx.InstallPath = ".go"
	}

	if !filepath.IsAbs(ctx.InstallPath) {
		rawPath := filepath.Join(xdg.Home, ctx.InstallPath)
		path, err := filepath.Abs(rawPath)
		if err != nil {
			fmt.Printf("Error computing absolute path for %s: %s\n", rawPath, err)
			os.Exit(1)
		}

		ctx.InstallPath = path
	}

	return ctx
}
