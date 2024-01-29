package deploy

import (
	"sync"
)

var (
	once            sync.Once
	defaultRegistry *Registry
)

func init() {
	once.Do(func() {
		defaultRegistry = NewRegistry()
	})
}

func RegisterResource(resourceType Namespace, name string, resource ResourceGenerator) {
	if err := defaultRegistry.Add(resourceType, name, resource); err != nil {
		panic(err)
	}
}

func DefaultRegistry() *Registry {
	return defaultRegistry
}
