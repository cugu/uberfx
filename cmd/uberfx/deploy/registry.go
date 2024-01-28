package deploy

import (
	"fmt"
)

type Registry struct {
	ResourceGenerators map[Namespace]map[string]ResourceGenerator
}

func NewRegistry() *Registry {
	return &Registry{
		ResourceGenerators: map[Namespace]map[string]ResourceGenerator{},
	}
}

func (r *Registry) Add(namespace Namespace, rsrcType string, rsrc ResourceGenerator) error {
	if _, ok := r.ResourceGenerators[namespace]; !ok {
		r.ResourceGenerators[namespace] = map[string]ResourceGenerator{}
	}

	if _, ok := r.ResourceGenerators[namespace][rsrcType]; ok {
		return fmt.Errorf("resource type %s/%s already registered", namespace, rsrcType)
	}

	r.ResourceGenerators[namespace][rsrcType] = rsrc

	return nil
}

func (r *Registry) ResourceGenerator(namespace Namespace, rsrcType string) (ResourceGenerator, error) {
	if rNamespace, ok := r.ResourceGenerators[namespace]; ok {
		if _, ok := rNamespace[rsrcType]; ok {
			return rNamespace[rsrcType], nil
		}
	}

	return nil, fmt.Errorf("resource type %s/%s not found", namespace, rsrcType)
}
