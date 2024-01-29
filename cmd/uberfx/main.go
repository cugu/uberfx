package main

import (
	_ "embed"
	"log"

	"github.com/alecthomas/kong"

	"github.com/cugu/uberfx/cmd/uberfx/deploy"
	"github.com/cugu/uberfx/cmd/uberfx/initialize"
	_ "github.com/cugu/uberfx/plugins"
)

type CLI struct {
	Verbose bool `help:"Enable verbose mode." short:"v" default:"false"`

	Init    initialize.Cmd `cmd:"" help:"Initialize a new project."`
	Deploy  deploy.Cmd     `cmd:"" help:"Deploy a project."`
	Install InstallCmd     `cmd:"" help:"Install a project."`
}

func main() {
	log.SetFlags(log.Lshortfile)

	var cli CLI
	ctx := kong.Parse(&cli)

	ctx.FatalIfErrorf(ctx.Run())
}
