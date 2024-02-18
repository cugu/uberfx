package initialize

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
)

var (
	uberspaceServerSubdomainRegex = regexp.MustCompile(`^[a-z]+$`)
	uberspaceServerRegex          = regexp.MustCompile(`^[a-z]+\.uberspace\.de$`)
	uberspaceUsernameRegex        = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	domainRegex                   = regexp.MustCompile(`^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`)
)

func interactiveInput(c *Cmd) error {
	// if we already have all the information we need, we don't need to prompt
	if c.Name != "" && c.GoVersion != "" && c.ModulePath != "" && c.UberspaceServer != "" && c.UberspaceUser != "" && c.Domain != "" {
		return nil
	}

	var style = lipgloss.NewStyle().
		Bold(true).
		Padding(1).
		MaxWidth(80)

	fmt.Println(style.Render("Create a new uberfx project"))

	if err := huh.NewForm(huh.NewGroup(
		nameField(&c.Name),
		goVersionField(&c.GoVersion),
		modulePathField(&c.ModulePath),
		uberspaceServerField(&c.UberspaceServer),
		uberspaceUserField(&c.UberspaceUser),
		domainField(&c.Domain),
	)).Run(); err != nil {
		return err
	}

	if !strings.HasSuffix(c.UberspaceServer, ".uberspace.de") {
		c.UberspaceServer += ".uberspace.de"
	}

	return nil
}

func nameField(name *string) huh.Field {
	return huh.NewInput().
		Title("Choose a name for your project").
		Placeholder("MyApp").
		Value(name).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("name cannot be empty")
			}

			return nil
		})
}

func goVersionField(goVersion *string) huh.Field {
	return huh.NewInput().
		Title("Choose a Go version (e.g. 1.22.0)").
		Placeholder("1.22.0").
		Value(goVersion).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("go version cannot be empty")
			}

			return nil
		})
}

func modulePathField(modulePath *string) huh.Field {
	return huh.NewInput().
		Title("Choose a Go module path (e.g. github.com/any/repository)").
		Placeholder("github.com/any/repository").
		Value(modulePath).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("package name cannot be empty")
			}

			return nil
		})
}

func uberspaceServerField(server *string) huh.Field {
	return huh.NewInput().
		Title("What is your Uberspace server? (e.g. stardust or stardust.uberspace.de)").
		Placeholder("stardust or stardust.uberspace.de").
		Value(server).
		Validate(func(s string) error {
			if !uberspaceServerSubdomainRegex.MatchString(s) && !uberspaceServerRegex.MatchString(s) {
				return errors.New("invalid Uberspace server name")
			}

			return nil
		})
}

func uberspaceUserField(user *string) huh.Field {
	return huh.NewInput().
		Title("What is your Uberspace user?").
		Placeholder("isabell").
		Value(user).
		Validate(func(s string) error {
			if !uberspaceUsernameRegex.MatchString(s) {
				return errors.New("invalid Uberspace user name")
			}

			return nil
		})
}

func domainField(domain *string) huh.Field {
	return huh.NewInput().
		Title("What is your domain?").
		Placeholder("www.yourname.uber.space").
		Value(domain).
		Validate(func(s string) error {
			if !domainRegex.MatchString(s) {
				return errors.New("invalid domain")
			}

			return nil
		})
}
