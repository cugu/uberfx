package plugins

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/cugu/uberfx/cmd/uberfx/deploy"
	"github.com/cugu/uberfx/ssh"
)

func init() {
	deploy.RegisterResource(deploy.NamespaceDeploy, "uberspace", NewUberspaceDeployer)
}

type UberspaceDeployer struct {
	In struct {
		Source string `json:"source" hcl:"source"`

		Username string `json:"username" hcl:"username"`
		Password string `json:"-" hcl:"password"`
		Address  string `json:"address" hcl:"address"`
		Port     int    `json:"port" hcl:"port"`

		Domain string            `json:"domain" hcl:"domain"`
		Env    map[string]string `json:"env" hcl:"env,optional"`
	} `json:"input"`
}

func serviceIni(u *UberspaceDeployer, binaryPath string) string {
	var envs []string

	if len(u.In.Env) > 0 {
		for k, v := range u.In.Env {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}
	}

	env := ""
	if len(envs) > 0 {
		env = "environment=" + strings.Join(envs, ",") + "\n"
	}

	iniTemplate := "[program:%s]\ncommand=%s :%d\nstartsecs=60\n%s"

	return fmt.Sprintf(iniTemplate, u.In.Domain, binaryPath, u.In.Port, env)
}

func NewUberspaceDeployer(body hcl.Body, ectx *hcl.EvalContext) (deploy.Resource, error) {
	u := &UberspaceDeployer{}
	if diags := gohcl.DecodeBody(body, ectx, &u.In); diags.HasErrors() {
		return nil, diags
	}

	randomString := randomString(5)
	binaryPath := fmt.Sprintf("/home/%s/bin/%s-%s", u.In.Username, u.In.Domain, randomString)
	iniPath := fmt.Sprintf("/home/%s/etc/services.d/%s.ini", u.In.Username, u.In.Domain)

	r, err := ssh.NewRemote(u.In.Address, u.In.Username, u.In.Password)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	slog.InfoContext(ctx, "copy app")
	if err := r.Copy(ctx, u.In.Source, binaryPath); err != nil {
		return nil, fmt.Errorf("copying server: %w", err)
	}

	slog.InfoContext(ctx, "chmod server")
	if err := r.Run(fmt.Sprintf("chmod +x %s", binaryPath)); err != nil {
		return nil, fmt.Errorf("chmod server: %w", err)
	}

	slog.InfoContext(ctx, "uberspace web domain add")
	_ = r.Run(fmt.Sprintf("uberspace web domain add %s", u.In.Domain))

	slog.InfoContext(ctx, "uberspace web backend set")
	if err := r.Run(fmt.Sprintf("uberspace web backend set %s --http --port %d", u.In.Domain, u.In.Port)); err != nil {
		return nil, fmt.Errorf("setting web backend: %w", err)
	}

	slog.InfoContext(ctx, "write services.d ini file")
	if err := r.Write(ctx, iniPath, serviceIni(u, binaryPath)); err != nil {
		return nil, fmt.Errorf("writing services.d: %w", err)
	}

	slog.Debug("Reloading service")
	if err := r.Run("supervisorctl reread"); err != nil {
		return nil, fmt.Errorf("supervisorctl reread error: %w", err)
	}

	slog.Debug("Updating service")
	if err := r.Run("supervisorctl update"); err != nil {
		return nil, fmt.Errorf("supervisorctl update error: %w", err)
	}

	slog.InfoContext(ctx, "supervisorctl start")
	if err := r.Run(fmt.Sprintf("supervisorctl start %s", u.In.Domain)); err != nil {
		slog.InfoContext(ctx, "supervisorctl restart")
		if err := r.Run(fmt.Sprintf("supervisorctl restart %s", u.In.Domain)); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func randomString(l int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, l)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
