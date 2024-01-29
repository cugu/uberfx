package plugins

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"

	deploy "github.com/cugu/uberfx/cmd/uberfx/deploy"
)

func init() {
	deploy.RegisterResource(deploy.NamespaceVar, "string", NewStringVar)
}

type StringVar struct {
	In struct {
		Name  string `json:"name" hcl:"name"`
		Value string `json:"value,omitempty" hcl:"value,optional"`
	} `json:"input"`

	Out struct {
		Source string `json:"source"`
		Value  string `json:"value"`
	} `json:"output"`
}

func NewStringVar(body hcl.Body, ctx *hcl.EvalContext) (deploy.Resource, error) {
	v := &StringVar{}
	if diags := gohcl.DecodeBody(body, ctx, &v.In); diags.HasErrors() {
		return nil, diags
	}

	if flags, ok := ctx.Variables["flags"]; ok {
		if flags.Type().HasAttribute(v.In.Name) {
			attr := flags.GetAttr(v.In.Name)

			v.Out.Source = "flag"
			v.Out.Value = attr.AsString()

			return v, nil
		}
	}

	if envValue, ok := os.LookupEnv("UBERFX_VAR_" + v.In.Name); ok {
		v.Out.Source = "env"
		v.Out.Value = envValue

		return v, nil
	}

	if v.In.Value != "" {
		v.Out.Source = "hcl"
		v.Out.Value = v.In.Value

		return v, nil
	}

	return v, nil
}

func (w *StringVar) Output() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"source": cty.StringVal(w.Out.Source),
		"value":  cty.StringVal(w.Out.Value),
	})
}
