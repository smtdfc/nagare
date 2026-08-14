package tool_declaration

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	tool_logger "github.com/smtdfc/nagare/core/tool/logger"
)

type ToolCallback[I any, O any] func(ctx domains.Context, args I) (O, error)

type ToolDeclaration[I any, O any] struct {
	Name        string
	Args        I
	Description string
	Callback    ToolCallback[I, O]
}

func (t *ToolDeclaration[I, O]) GetName() string {
	return t.Name
}

func (t *ToolDeclaration[I, O]) GetArgs() string {
	var r I
	reflector := &jsonschema.Reflector{
		DoNotReference: true,
	}
	schema := reflector.Reflect(&r)

	data, err := json.Marshal(schema)
	if err != nil {
		tool_logger.ToolLogger.Error("Failed to marshal tool schema", "error", err)
		return "{}"
	}

	return string(data)
}

func (t *ToolDeclaration[I, O]) GetDescription() string {
	return t.Description
}

func (t *ToolDeclaration[I, O]) Execute(ctx domains.Context, argsJson string) (string, error) {
	var args I
	err := json.Unmarshal([]byte(argsJson), &args)
	if err != nil {
		tool_logger.ToolLogger.Error("Failed to unmarshal tool arguments", "error", err)
		return "{}", custom_errors.NewToolError("Tool validation failed: One or more arguments provided to the tool are incorrect or improperly formatted.", t.Name)
	}

	result, err := t.Callback(ctx, args)
	if err != nil {
		return "{}", err
	}

	resultJson, err := json.Marshal(&result)
	if err != nil {
		tool_logger.ToolLogger.Error("Failed to marshal tool result", "error", err)
		return "{}", custom_errors.NewToolError("Tool processing error: An unexpected error occurred while processing the tool's output.", t.Name)
	}

	return string(resultJson), nil
}

func Declare[I any, O any](name string, description string, cb ToolCallback[I, O]) *ToolDeclaration[I, O] {
	return &ToolDeclaration[I, O]{
		Name:        name,
		Description: description,
		Callback:    cb,
	}
}
