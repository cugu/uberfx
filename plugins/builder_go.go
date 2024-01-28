package plugins

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"

	"github.com/cugu/uberfx/cmd/uberfx/deploy"
)

func init() {
	deploy.RegisterResource(deploy.NamespaceBuild, "go", NewGoBuilder)
}

type GoBuilder struct {
	In struct {
		Path string `json:"path" hcl:"path"`
	} `json:"input"`

	Out struct {
		Path string `json:"path"`
	} `json:"output"`
}

func NewGoBuilder(body hcl.Body, ectx *hcl.EvalContext) (deploy.Resource, error) {
	w := &GoBuilder{}
	if diags := gohcl.DecodeBody(body, ectx, &w.In); diags.HasErrors() {
		return nil, diags
	}

	ctx := context.Background()

	slog.InfoContext(ctx, "creating build directory")
	buildDir, err := os.MkdirTemp("", "uberfx")
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "creating cache directory")
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "looking for go")
	goPath, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}

	w.Out.Path = fmt.Sprintf("%s/main.wasm", buildDir)

	slog.InfoContext(ctx, fmt.Sprintf("go build -o %s %s", w.Out.Path, w.In.Path))
	goCmd := exec.Cmd{
		Path: goPath,
		Args: []string{"go", "build", "-o", w.Out.Path, w.In.Path},
		Env: []string{
			"GOOS=wasip1",
			"GOARCH=wasm",
			fmt.Sprintf("GOCACHE=%s/uberfx/gocache", cacheDir),
			fmt.Sprintf("GOMODCACHE=%s/uberfx/gomodcache", cacheDir),
		},
	}

	out, err := goCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error building: %w: %s", err, out)
	}

	return w, nil
}

func (w *GoBuilder) Output() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"output": cty.StringVal(w.Out.Path),
	})
}
