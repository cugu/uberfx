package deploy

import (
	"encoding/json"
)

type Namespace int

const (
	NamespaceVar Namespace = iota
	NamespaceProvider
	NamespaceService
	NamespaceBuild
	NamespaceDeploy
)

var ResourceTypes = []Namespace{
	NamespaceVar,
	NamespaceProvider,
	NamespaceService,
	NamespaceBuild,
	NamespaceDeploy,
}

func (r Namespace) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r Namespace) String() string {
	switch r {
	case NamespaceVar:
		return "var"
	case NamespaceProvider:
		return "provider"
	case NamespaceService:
		return "service"
	case NamespaceBuild:
		return "build"
	case NamespaceDeploy:
		return "deploy"
	default:
		panic("unknown namespace")
	}
}
