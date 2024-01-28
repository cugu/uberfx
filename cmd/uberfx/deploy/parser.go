package deploy

import (
	"fmt"
	"maps"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type Parser struct {
	registry *Registry

	allFlags map[string]string

	State *State `json:"state"`
}

type State struct {
	Resources []*ResourceEntry `json:"resources"`
}

func NewParser(allFlags map[string]string) *Parser {
	return &Parser{
		registry: DefaultRegistry(),
		allFlags: allFlags,
		State:    &State{},
	}
}

func (p *Parser) Decode(filename string, data []byte) error {
	file, diags := hclsyntax.ParseConfig(data, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}

	bodyContent, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: NamespaceVar.String(), LabelNames: []string{"type", "name"}},
			{Type: NamespaceProvider.String(), LabelNames: []string{"type", "name"}},
			{Type: NamespaceService.String(), LabelNames: []string{"type", "name"}},
			{Type: NamespaceBuild.String(), LabelNames: []string{"type", "name"}},
			{Type: NamespaceDeploy.String(), LabelNames: []string{"type", "name"}},
		},
	})
	if diags.HasErrors() {
		return diags
	}

	if err := p.parseResources(bodyContent); err != nil {
		return fmt.Errorf("error parsing resources: %w", err)
	}

	return nil
}

func (p *Parser) parseResources(bodyContent *hcl.BodyContent) error {
	for _, res := range ResourceTypes {
		resourceBlocks := bodyContent.Blocks.OfType(res.String())
		for _, resourceBlock := range resourceBlocks {
			if len(resourceBlock.Labels) != 2 {
				return fmt.Errorf("resource %s must have exactly two labels", resourceBlock.Type)
			}

			id := ResourceID{res, resourceBlock.Labels[0], resourceBlock.Labels[1]}

			resourceGenerator, err := p.registry.ResourceGenerator(res, resourceBlock.Labels[0])
			if err != nil {
				return err
			}

			rsrc, err := resourceGenerator(resourceBlock.Body, p.evalContext())
			if err != nil {
				return fmt.Errorf("error generating resource %s: %w", id.String(), err)
			}

			p.State.Resources = append(p.State.Resources, &ResourceEntry{
				ID:       id,
				Resource: rsrc,
			})
		}
	}

	return nil
}

func (p *Parser) evalContext() *hcl.EvalContext {
	all := map[string]cty.Value{}
	for k, v := range p.allFlags {
		all[k] = cty.StringVal(v)
	}

	variables := map[string]cty.Value{
		"flags": cty.ObjectVal(all),
	}

	maps.Copy(variables, p.resourceContext())

	return &hcl.EvalContext{
		Variables: variables,
	}
}

func (p *Parser) resourceContext() map[string]cty.Value {
	namespaces := map[Namespace]map[string]map[string]cty.Value{}

	for _, s := range p.State.Resources {
		if outputer, ok := s.Resource.(Outputer); ok {
			if _, ok := namespaces[s.ID.Namespace]; !ok {
				namespaces[s.ID.Namespace] = map[string]map[string]cty.Value{}
			}

			if _, ok := namespaces[s.ID.Namespace][s.ID.Type]; !ok {
				namespaces[s.ID.Namespace][s.ID.Type] = map[string]cty.Value{}
			}

			namespaces[s.ID.Namespace][s.ID.Type][s.ID.Name] = outputer.Output()
		}
	}

	variables := map[string]cty.Value{}

	for namespace, namespaceResources := range namespaces {
		res := map[string]cty.Value{}

		for name, ids := range namespaceResources {
			res[name] = cty.ObjectVal(ids)
		}

		variables[namespace.String()] = cty.ObjectVal(res)
	}

	return variables
}
