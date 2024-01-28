package deploy

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type ResourceID struct {
	Namespace Namespace `json:"namespace"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
}

func (r ResourceID) String() string {
	return fmt.Sprintf("%s.%s.%s", r.Namespace, r.Type, r.Name)
}

func (r ResourceID) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

type ResourceEntry struct {
	ID       ResourceID `json:"id"`
	Resource Resource   `json:"resource"`
}

type Resource interface{}

type Outputer interface {
	Output() cty.Value
}

type ResourceGenerator func(body hcl.Body, ctx *hcl.EvalContext) (Resource, error)
