package plugins

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"

	"github.com/cugu/uberfx/cmd/uberfx/deploy"
	"github.com/cugu/uberfx/ssh"
)

func init() {
	deploy.RegisterResource(deploy.NamespaceService, "uberspace_mysql_db", NewUberspaceMySQLDatabase)
}

type UberspaceMySQLDatabase struct {
	In struct {
		Username string `json:"username" hcl:"username"`
		Password string `json:"-" hcl:"password"`
		Address  string `json:"address" hcl:"address"`
		Suffix   string `json:"suffix" hcl:"suffix"`
	} `json:"input"`

	Out struct {
		Name string `json:"name" hcl:"name"`
	} `json:"output"`
}

func NewUberspaceMySQLDatabase(body hcl.Body, ectx *hcl.EvalContext) (deploy.Resource, error) {
	u := &UberspaceMySQLDatabase{}
	if diags := gohcl.DecodeBody(body, ectx, &u.In); diags.HasErrors() {
		return nil, diags
	}

	r, err := ssh.NewRemote(u.In.Address, u.In.Username, u.In.Password)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s_%s", u.In.Username, u.In.Suffix)

	_, _ = r.Run(fmt.Sprintf(`mysql -e "CREATE DATABASE %s"`, name))

	u.Out.Name = name

	return u, nil
}

func (w *UberspaceMySQLDatabase) Output() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"name": cty.StringVal(w.Out.Name),
	})
}
