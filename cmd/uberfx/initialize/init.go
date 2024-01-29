package initialize

import (
	_ "embed"
	"os"
)

//go:embed template/uberfx.hcl
var ConfigTemplate string

//go:embed template/main.go
var MainTemplate string

type Cmd struct {
	Dir string `help:"Directory to initialize." arg:"" type:"path"`
}

func (c *Cmd) Run() error {
	if err := os.Mkdir(c.Dir, 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(c.Dir+"/main.go", []byte(MainTemplate), 0o644); err != nil {
		return err
	}

	if err := os.WriteFile(c.Dir+"/uberfx.hcl", []byte(ConfigTemplate), 0o644); err != nil {
		return err
	}

	return nil
}
