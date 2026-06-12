package tool

import "github.com/invopop/jsonschema"

type ParameterDefinition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func ExtractParameters(args any) ParameterDefinition {
	schema := jsonschema.Reflect(args)

	var defName string
	for k := range schema.Definitions {
		defName = k
		break
	}

	rootDef := schema.Definitions[defName]

	props := make(map[string]Property)

	for pair := rootDef.Properties.Oldest(); pair != nil; pair = pair.Next() {
		props[pair.Key] = Property{
			Type:        pair.Value.Type,
			Description: pair.Value.Description,
		}
	}

	return ParameterDefinition{
		Type:       "object",
		Properties: props,
		Required:   rootDef.Required,
	}
}
