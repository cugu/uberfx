package plugins

import (
	"context"
	"fmt"
	"log/slog"

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

		Domain string            `json:"domain" hcl:"domain"`
		Env    map[string]string `json:"env" hcl:"env,optional"`
	} `json:"input"`
}

func NewUberspaceDeployer(body hcl.Body, ectx *hcl.EvalContext) (deploy.Resource, error) {
	u := &UberspaceDeployer{}
	if diags := gohcl.DecodeBody(body, ectx, &u.In); diags.HasErrors() {
		return nil, diags
	}

	r, err := ssh.NewRemote(u.In.Address, u.In.Username, u.In.Password)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	slog.InfoContext(ctx, "scp wasm server")
	if err := r.Copy(ctx, u.In.Source, fmt.Sprintf("/home/%s/bin/wasi/%s.wasm", u.In.Username, u.In.Domain)); err != nil {
		return nil, fmt.Errorf("copying wasm server: %w", err)
	}

	if len(u.In.Env) > 0 {
		slog.InfoContext(ctx, "scp wasm env")
		envContent := ""
		for k, v := range u.In.Env {
			envContent += fmt.Sprintf("%s=%s\n", k, v)
		}
		if err := r.Write(ctx, fmt.Sprintf("/home/%s/bin/wasi/%s.wasm.env", u.In.Username, u.In.Domain), envContent); err != nil {
			return nil, fmt.Errorf("writing wasm env: %w", err)
		}
	}

	slog.InfoContext(ctx, "uberspace web domain add")
	_ = r.Run(fmt.Sprintf("uberspace web domain add %s", u.In.Domain))

	slog.InfoContext(ctx, "uberspace web backend set")
	if err := r.Run(fmt.Sprintf("uberspace web backend set %s --http --port 8080", u.In.Domain)); err != nil {
		return nil, fmt.Errorf("setting web backend: %w", err)
	}

	slog.InfoContext(ctx, "supervisorctl restart uberfx-server")
	if err := r.Run("supervisorctl restart uberfx-server"); err != nil {
		return nil, fmt.Errorf("setting web backend: %w", err)
	}

	return u, nil
}
