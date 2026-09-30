package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/invopop/jsonschema"
	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

type PluginTool interface {
	GetName() string
	GetDescription() string
	GetArgsSchema() string
	GetCategories() []string
	Execute(ctx *context.Context, args string) (string, error)
}

type ListPluginTool []PluginTool

type BaseToolCallback[I any, O any] func(ctx *context.Context, args *I) (*O, error)
type BaseTool[I any, O any] struct {
	Name        string
	Description string
	Callback    BaseToolCallback[I, O]
	Categories  []string
}

// Execute implements [Tool].
func (b *BaseTool[I, O]) Execute(ctx *context.Context, argRaw string) (string, error) {
	args, err := helpers.UnmarshalJson[I](argRaw)
	if err != nil {
		return "{}", ErrIncorrectToolArgs
	}

	result, err := b.Callback(ctx, args)
	if err != nil {
		return "{}", err
	}

	resultJson, err := helpers.MarshalJson(&result)
	if err != nil {
		return "{}", ErrMarshalToolResultFailed
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

// GetCategories implements [Tool].
func (b *BaseTool[I, O]) GetCategories() []string {
	return b.Categories
}

func DefineTool[I any, O any](
	name string,
	description string,
	cb BaseToolCallback[I, O],
	categories []string,
) PluginTool {
	tool := &BaseTool[I, O]{
		Name:        name,
		Description: description,
		Callback:    cb,
		Categories:  categories,
	}

	return tool
}

func (p *PluginClient) RegisterPluginTool(ctx context.Context, tool PluginTool) error {
	payload := plugin_dtos.RegisterPluginToolEventPayload{
		Tool: plugin_dtos.PluginTool{
			Name:        tool.GetName(),
			Description: tool.GetDescription(),
			Args:        tool.GetArgsSchema(),
			Categories:  tool.GetCategories(),
		},
	}

	resp, err := sendAndWait[plugin_dtos.RegisterPluginToolSuccessEventPayload, plugin_dtos.RegisterPluginToolFailedEventPayload](
		p,
		ctx,
		plugin_dtos.RegisterPluginToolEvent,
		payload,
		plugin_dtos.RegisterPluginToolFailedEvent,
		plugin_dtos.RegisterPluginToolSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to register plugin tool", "error", err, "tool", tool.GetName())
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to register plugin tool", "error", resp.Error, "tool", tool.GetName())
		return errors.New(resp.Error.Cause)
	}

	p.mu.Lock()
	p.tools = append(p.tools, tool)
	p.mu.Unlock()

	return nil
}

func (p *PluginClient) RegisterToolCategories(ctx context.Context, categories map[string]string) error {
	payload := plugin_dtos.RegisterToolCategoriesEventPayload{
		Categories: categories,
	}

	resp, err := sendAndWait[plugin_dtos.RegisterToolCategoriesSuccessEventPayload, plugin_dtos.RegisterToolCategoriesFailedEventPayload](
		p,
		ctx,
		plugin_dtos.RegisterToolCategoriesEvent,
		payload,
		plugin_dtos.RegisterToolCategoriesFailedEvent,
		plugin_dtos.RegisterToolCategoriesSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to register tool categories", "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to register tool categories", "error", resp.Error)
		return errors.New(resp.Error.Cause)
	}

	return nil
}
