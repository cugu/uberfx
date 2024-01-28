package deploy

import (
	"encoding/json"
	"log"
	"os"
)

type Cmd struct {
	Var    map[string]string `help:"Pass extra variables."`
	Config string            `help:"Config file." type:"path" default:"uberfx.hcl"`
}

func (c *Cmd) Run() error {
	_, err := os.Stat(c.Config)
	if err != nil {
		return err
	}

	log.Println("deploying", c.Config)

	config, err := os.ReadFile(c.Config)
	if err != nil {
		return err
	}

	d := NewParser(c.Var)

	if err := d.Decode(c.Config, config); err != nil {
		return err
	}

	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	log.Println(string(b))

	return nil
}
