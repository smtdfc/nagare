package plugin

import (
	"errors"
	"time"

	"github.com/google/uuid"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/tool"
)

type Tool struct {
	Name        string
	Args        string
	Description string
	Categories  []string
	PluginID    string
	eventBus    *event_bus.CoreEventBus
}

// Execute implements [tool.Tool].
func (p *Tool) Execute(ctx *core_context.ExecuteContext, args string) (string, error) {
	requestID := uuid.NewString()
	resultChannel, unsubscribe := p.eventBus.Subscribe(event_bus.PluginToolCallResultEvent)
	defer unsubscribe()

	p.eventBus.Publish(ctx, event_bus.PluginToolCallEvent, &event_bus.PluginToolCallEventPayload{
		RequestID: requestID,
		PluginID:  p.PluginID,
		Name:      p.Name,
		Args:      args,
	})

	timer := time.NewTimer(60 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return "{}", ctx.Err()
		case <-timer.C:
			return "{}", errors.New("plugin tool call timed out")
		case payload, ok := <-resultChannel:
			if !ok {
				return "{}", errors.New("plugin tool result subscription closed")
			}
			result, ok := payload.(*event_bus.PluginToolCallResultEventPayload)
			if !ok || result.RequestID != requestID {
				continue
			}
			if result.Error != "" {
				return "{}", errors.New(result.Error)
			}
			return result.Result, nil
		}
	}
}

// GetArgsSchema implements [tool.Tool].
func (p *Tool) GetArgsSchema() string {
	return p.Args
}

// GetBindings implements [tool.Tool].
func (p *Tool) GetBindings() tool.Bindings {
	return nil
}

// GetCategories implements [tool.Tool].
func (p *Tool) GetCategories() []string {
	return p.Categories
}

// GetDescription implements [tool.Tool].
func (p *Tool) GetDescription() string {
	return p.Description
}

// GetName implements [tool.Tool].
func (p *Tool) GetName() string {
	return p.Name
}

// WithBindings implements [tool.Tool].
func (p *Tool) WithBindings(tool.Bindings) tool.Tool {
	return p
}

func NewPluginTool(name string, args string, description string, categories []string, pluginID string, eventBus *event_bus.CoreEventBus) *Tool {
	return &Tool{
		Name:        name,
		Args:        args,
		Description: description,
		Categories:  categories,
		PluginID:    pluginID,
		eventBus:    eventBus,
	}
}
