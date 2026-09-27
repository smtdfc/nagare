package client

import (
	"context"
	"encoding/json"
	"errors"
	"uuid"

	"github.com/invopop/jsonschema"
	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/shared/helpers"
)

type DynamicTool interface {
	GetName() string
	GetDescription() string
	GetArgsSchema() string
	Execute(ctx *context.Context, args string) (string, error)
}

type ListDynamicTool []DynamicTool

type BaseToolCallback[I any, O any] func(ctx *context.Context, args *I) (*O, error)
type BaseTool[I any, O any] struct {
	Name        string
	Description string
	Callback    BaseToolCallback[I, O]
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

func DefineTool[I any, O any](
	name string,
	description string,
	cb BaseToolCallback[I, O],
) DynamicTool {
	tool := &BaseTool[I, O]{
		Name:        name,
		Description: description,
		Callback:    cb,
	}

	return tool
}

func (p *PluginClient) UseDynamicTool(tool DynamicTool) {
	p.dynamicTools = append(p.dynamicTools, tool)
}

func (p *PluginClient) RegisterDynamicTools(ctx context.Context) error {
	requestID := uuid.NewV4().String()
	respChan := make(chan *websocket.Payload[any], 1)

	p.mu.Lock()
	p.pendingRequests[requestID] = respChan
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.pendingRequests, requestID)
		p.mu.Unlock()
	}()

	dtos := make([]plugin_dtos.DynamicTool, len(p.dynamicTools))
	for i, tool := range p.dynamicTools {
		dtos[i] = plugin_dtos.DynamicTool{
			Name:        tool.GetName(),
			Description: tool.GetDescription(),
			ArgsSchema:  tool.GetArgsSchema(),
		}
	}
	err := p.connector.Send(
		plugin_dtos.RegisterDynamicToolsEvent,
		plugin_dtos.RegisterDynamicToolsPayload{
			Tools: dtos,
		},
		requestID,
	)
	if err != nil {
		p.Logger.Error("failed to register dynamic tools", "error", err)
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case resp := <-respChan:
		if resp.Event == plugin_dtos.RegisterDynamicToolsErrorEvent {
			payload, _ := GetData[plugin_dtos.RegisterDynamicToolsErrorPayload](resp)
			p.Logger.Error("failed to register dynamic tools", "error", payload)
			return errors.New(payload.Cause)
		}

		return nil
	}
}
