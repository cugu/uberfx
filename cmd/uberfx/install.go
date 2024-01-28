package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/cugu/uberfx/ssh"
)

const servicesD = `[program:uberfx-server]
command=/home/%s/bin/uberfx-server --address :8080 --wasi-dir /home/%s/bin/wasi --debug
startsecs=60
`

type InstallCmd struct {
	Version  string `help:"Version of uberfx-server to install." type:"string" default:"latest"`
	Address  string `help:"Address of the server." type:"string" required:""`
	Username string `help:"Username of the server." type:"string" required:""`
	Password string `help:"Password of the server." type:"string" required:""`
}

func (c *InstallCmd) Run() error {
	ctx := context.Background()

	r, err := ssh.NewRemote(c.Address, c.Username, c.Password)
	if err != nil {
		return err
	}

	downloadUrl, err := downloadUrl(ctx, c.Version)
	if err != nil {
		return fmt.Errorf("getting download url: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("curl -L %s -o uberfx-server.tar.gz", downloadUrl))
	if err := r.Run(fmt.Sprintf("curl -L %s -o uberfx-server.tar.gz", downloadUrl)); err != nil {
		return fmt.Errorf("copying uberfx-server: %w", err)
	}

	slog.InfoContext(ctx, "supervisorctl stop uberfx-server")
	if err := r.Run("supervisorctl stop uberfx-server"); err != nil {
		return fmt.Errorf("stop uberfx-server: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("tar -xzf uberfx-server.tar.gz -C /home/%s/bin", c.Username))
	if err := r.Run(fmt.Sprintf("tar -xzf uberfx-server.tar.gz -C /home/%s/bin", c.Username)); err != nil {
		return fmt.Errorf("extracting uberfx-server: %w", err)
	}

	slog.InfoContext(ctx, "chmod +x uberfx-server")
	if err := r.Run(fmt.Sprintf("chmod +x /home/%s/bin/uberfx-server", c.Username)); err != nil {
		return fmt.Errorf("chmod +x uberfx-server: %w", err)
	}

	slog.InfoContext(ctx, "supervisorctl reread")
	if err := r.Run("supervisorctl reread"); err != nil {
		return fmt.Errorf("rereading supervisorctl: %w", err)
	}

	slog.InfoContext(ctx, "supervisorctl update")
	if err := r.Run("supervisorctl update"); err != nil {
		return fmt.Errorf("updating supervisorctl: %w", err)
	}

	slog.InfoContext(ctx, "supervisorctl start uberfx-server")
	if err := r.Run("supervisorctl start uberfx-server"); err != nil {
		return fmt.Errorf("restarting uberfx-server: %w", err)
	}

	slog.InfoContext(ctx, "rm uberfx-server.tar.gz")
	if err := r.Run("rm uberfx-server.tar.gz"); err != nil {
		return fmt.Errorf("removing uberfx-server.tar.gz: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("mkdir -p /home/%s/bin/wasi", c.Username))
	if err := r.Run(fmt.Sprintf("mkdir -p /home/%s/bin/wasi", c.Username)); err != nil {
		return fmt.Errorf("mkdir -p /home/%s/bin/wasi: %w", c.Username, err)
	}

	slog.InfoContext(ctx, "create services.d/wasm-server.ini")
	content := fmt.Sprintf(servicesD, c.Username, c.Username)
	if err := r.Write(ctx, fmt.Sprintf("/home/%s/etc/services.d/wasm-server.ini", c.Username), content); err != nil {
		return fmt.Errorf("writing services.d: %w", err)
	}

	return nil
}

var versionRegex = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)

func downloadUrl(ctx context.Context, version string) (string, error) {
	ubuntuURL := "https://github.com/cugu/uberfx-server/releases/download/%s/uberfx-server_Linux_x86_64.tar.gz"

	if versionRegex.MatchString(version) {
		return fmt.Sprintf(ubuntuURL, version), nil
	}

	if version != "latest" {
		return "", errors.New("invalid version")
	}

	url := "https://api.github.com/repos/cugu/uberfx-server/releases/latest"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %s", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return fmt.Sprintf(ubuntuURL, release.TagName), nil
}
