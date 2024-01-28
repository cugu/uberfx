package plugins

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
	ini "gopkg.in/ini.v1"

	"github.com/cugu/uberfx/cmd/uberfx/deploy"
	"github.com/cugu/uberfx/ssh"
)

func init() {
	deploy.RegisterResource(deploy.NamespaceService, "uberspace_mysql", NewUberspaceMySQLService)
}

type UberspaceMySQLService struct {
	In struct {
		Username string `json:"username" hcl:"username"`
		Password string `json:"-" hcl:"password"`
		Address  string `json:"address" hcl:"address"`
	} `json:"input"`

	Out struct {
		User       string `json:"user" hcl:"user"`
		Password   string `json:"password" hcl:"password"`
		UserRO     string `json:"user_ro" hcl:"user_ro"`
		PasswordRO string `json:"password_ro" hcl:"password_ro"`
	} `json:"output"`
}

func NewUberspaceMySQLService(body hcl.Body, ectx *hcl.EvalContext) (deploy.Resource, error) {
	u := &UberspaceMySQLService{}
	if diags := gohcl.DecodeBody(body, ectx, &u.In); diags.HasErrors() {
		return nil, diags
	}

	r, err := ssh.NewRemote(u.In.Address, u.In.Username, u.In.Password)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	cnf, err := r.Read(ctx, fmt.Sprintf("/home/%s/.my.cnf", u.In.Username))
	if err != nil {
		return nil, fmt.Errorf("reading .my.cnf: %w", err)
	}

	cfg, err := ini.Load(cnf)
	if err != nil {
		return nil, fmt.Errorf("loading .my.cnf: %w", err)
	}

	section, err := cfg.GetSection("client")
	if err != nil {
		return nil, fmt.Errorf("getting client section: %w", err)
	}

	u.Out.User = section.Key("user").String()
	u.Out.Password = section.Key("password").String()

	roSection, err := cfg.GetSection("clientreadonly")
	if err != nil {
		return nil, fmt.Errorf("getting clientreadonly section: %w", err)
	}

	u.Out.UserRO = roSection.Key("user").String()
	u.Out.PasswordRO = roSection.Key("password").String()

	return u, nil
}

func (w *UberspaceMySQLService) Output() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"user":        cty.StringVal(w.Out.User),
		"password":    cty.StringVal(w.Out.Password),
		"user_ro":     cty.StringVal(w.Out.UserRO),
		"password_ro": cty.StringVal(w.Out.PasswordRO),
	})
}
