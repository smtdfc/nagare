package tool

import (
	"encoding/json"
	"fmt"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/exceptions"
	nagare_collections "github.com/smtdfc/nagare/shared/collections"
)

type ToolCallback[T any, R any] func(ctx domains.AgentContext, args T) (R, error)

type ToolDeclaration[T any, R any] struct {
	Name        string
	Description string
	Callback    ToolCallback[T, R]
	Arguments   string
	ToolType    domains.ToolType
	Categories  []string
}

// GetGroup implements [Tool].
func (d *ToolDeclaration[T, R]) HasCategory(names []string) bool {
	return nagare_collections.HasStringIntersection(names, d.Categories)
}

// GetDesc implements [Tool].
func (d *ToolDeclaration[T, R]) GetDesc() string {
	return d.Description
}

func (d *ToolDeclaration[T, R]) GetName() string {
	return d.Name
}

func (d *ToolDeclaration[T, R]) Execute(ctx domains.AgentContext, argsRaw string) (string, error) {
	var args T

	if err := json.Unmarshal([]byte(argsRaw), &args); err != nil {
		return "", exceptions.NewToolException(fmt.Sprintf("invalid json params for tool %s: %s", d.Name, err), d.Name)
	}

	result, err := d.Callback(ctx, args)
	if err != nil {
		return "", exceptions.NewToolException(fmt.Sprintf("failed to call tool %s: %s", d.Name, err), d.Name)
	}

	if strResult, ok := any(result).(string); ok {
		return strResult, nil
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", exceptions.NewToolException(fmt.Sprintf("cannot marshal result of tool %s: %s", d.Name, err), d.Name)
	}

	return string(jsonData), nil
}

// func (d *ToolDeclaration[T, R]) ExecuteWithRawResult(ctx context.Context, argsRaw string) (any, error) {
// 	var args T
// 	var zero R

// 	if err := json.Unmarshal([]byte(argsRaw), &args); err != nil {
// 		return zero, exceptions.NewToolException(fmt.Sprintf("invalid json params for tool %s: %s", d.Name, err), d.Name)
// 	}

// 	result, err := d.Callback(ctx, args)
// 	if err != nil {
// 		return zero, exceptions.NewToolException(fmt.Sprintf("failed to call tool %s: %s", d.Name, err), d.Name)
// 	}

// 	return result, nil
// }

func (d *ToolDeclaration[T, R]) GetArgumentsSchema() string {
	return d.Arguments
}

func (d *ToolDeclaration[T, R]) GetType() domains.ToolType {
	return d.ToolType
}

func DeclareTool[T any, R any](
	name string,
	description string,
	cb ToolCallback[T, R],
	toolType domains.ToolType,
	categories []string,
) domains.Tool {
	var args T
	paramsDef := ExtractParameters(args)
	schemaBytes, err := json.Marshal(paramsDef)
	if err != nil {
		panic(fmt.Sprintf("error marshaling schema: %v", err))
	}

	// fmt.Println(string(schemaBytes))
	return &ToolDeclaration[T, R]{
		Name:        name,
		Description: description,
		Callback:    cb,
		Arguments:   string(schemaBytes),
		ToolType:    toolType,
		Categories:  categories,
	}
}
