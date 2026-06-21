package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/smtdfc/nagare/core/exceptions"
)

type Tool interface {
	GetName() string
	Execute(context.Context, string) (string, error)
	GetArgumentsSchema() string
	GetDesc() string
}

type ToolCallback[T any, R any] func(ctx context.Context, args T) (R, error)
type ToolDeclaration[T any, R any] struct {
	Name        string
	Description string
	Callback    ToolCallback[T, R]
	Arguments   string
}

// GetDesc implements [Tool].
func (d *ToolDeclaration[T, R]) GetDesc() string {
	return d.Description
}

func (d *ToolDeclaration[T, R]) GetName() string {
	return d.Name
}

func (d *ToolDeclaration[T, R]) Execute(ctx context.Context, argsRaw string) (string, error) {
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

func (d *ToolDeclaration[T, R]) GetArgumentsSchema() string {
	return d.Arguments
}

func DeclareTool[T any, R any](
	name string,
	description string,
	cb ToolCallback[T, R],
) Tool {
	var args T
	paramsDef := ExtractParameters(args)
	schemaBytes, err := json.Marshal(paramsDef)
	if err != nil {
		panic(fmt.Sprintf("error marshaling schema: %v", err))
	}

	return &ToolDeclaration[T, R]{
		Name:        name,
		Description: description,
		Callback:    cb,
		Arguments:   string(schemaBytes),
	}
}

type ListTool []Tool
