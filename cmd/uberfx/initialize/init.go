package initialize

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

//go:embed template
var templateFS embed.FS

type Cmd struct {
	Name            string `help:"Name of the project"`
	GoVersion       string `help:"Go version (e.g. 1.22.0)" default:"1.22.0"`
	ModulePath      string `help:"Go module path (e.g. github.com/yourname/yourapp)"`
	UberspaceServer string `help:"Uberspace server (e.g. stardust or stardust.uberspace.de)"`
	UberspaceUser   string `help:"Uberspace user"`
	Domain          string `help:"Domain"`
}

func (c *Cmd) Run() error {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if err := interactiveInput(c); err != nil {
			return err
		}
	} else {
		var issues []string
		if c.Name == "" {
			issues = append(issues, "name cannot be empty")
		}

		if c.GoVersion == "" {
			issues = append(issues, "go version cannot be empty")
		}

		if c.ModulePath == "" {
			issues = append(issues, "go module path cannot be empty")
		}

		if c.UberspaceServer == "" {
			issues = append(issues, "uberspace server cannot be empty")
		}

		if c.UberspaceUser == "" {
			issues = append(issues, "uberspace user cannot be empty")
		}

		if c.Domain == "" {
			issues = append(issues, "domain cannot be empty")
		}

		if len(issues) > 0 {
			return fmt.Errorf(strings.Join(issues, ", "))
		}
	}

	if err := validDest(c.Name); err != nil {
		return fmt.Errorf("could not create project %s: %w", c.Name, err)
	}

	return writeTemplates(c.Name, c)
}

func validDest(dest string) error {
	stat, err := os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if !stat.IsDir() {
		return errors.New("a file with the same name already exists")
	}

	dir, err := os.ReadDir(dest)
	if err != nil {
		return err
	}

	if len(dir) > 0 {
		return errors.New("directory already exists and is not empty")
	}

	return nil
}
