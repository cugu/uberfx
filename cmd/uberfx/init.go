package main

import (
	_ "embed"
	"os"

	"github.com/cugu/uberfx/example"
)

type InitCmd struct {
	Dir string `help:"Directory to initialize." arg:"" type:"path"`
}

func (c *InitCmd) Run() error {
	if err := os.Mkdir(c.Dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(c.Dir+"/main.go", []byte(example.MainTemplate), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(c.Dir+"/uberfx.hcl", []byte(example.ConfigTemplate), 0644); err != nil {
		return err
	}

	return nil
}
