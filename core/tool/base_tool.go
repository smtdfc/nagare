package tool

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/shared/helpers"
)

type BaseToolCallback[I any, O any] func(ctx context.Context, args *I) (*O, error)
type BaseTool[I any, O any] struct {
	Name        string
	Description string
	Callback    BaseToolCallback[I, O]
}

// Execute implements [Tool].
func (b *BaseTool[I, O]) Execute(ctx context.Context, argRaw string) (string, error) {
	args, err := helpers.UnmarshalJson[I](argRaw)
	if err != nil {
		return "{}", custom_errors.ErrIncorrectToolArgs
	}

	result, err := b.Callback(ctx, args)
	if err != nil {
		return "{}", err
	}

	resultJson, err := helpers.MarshalJson(&result)
	if err != nil {
		return "{}", custom_errors.ErrMarshalToolResultFailed
	}

	return string(resultJson), nil
}

// GetArgsSchema implements [Tool].
func (b *BaseTool[I, O]) GetArgsSchema() string {
	var r I
	reflector := &jsonschema.Reflector{
		DoNotReference: true,
	}
	schema := reflector.Reflect(&r)

	data, err := json.Marshal(schema)
	if err != nil {
		return "{}"
	}

	return string(data)
}

// GetDescription implements [Tool].
func (b *BaseTool[I, O]) GetDescription() string {
	return b.Description
}

// GetName implements [Tool].
func (b *BaseTool[I, O]) GetName() string {
	return b.Name
}

func DefineTool[I any, O any](
	name string,
	description string,
	cb BaseToolCallback[I, O],
) Tool {
	return &BaseTool[I, O]{
		Name:        name,
		Description: description,
		Callback:    cb,
	}
}
