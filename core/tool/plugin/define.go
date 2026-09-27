package plugin

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type PluginTool struct {
	Name        string
	Args        string
	Description string
	Categories  []string
	PluginID    string
}

// Execute implements [tool.Tool].
func (p *PluginTool) Execute(ctx *context.ExecuteContext, args string) (string, error) {
	fmt.Println("Called", p.Name, "with args:", args)

	return "{}", nil
}

// GetArgsSchema implements [tool.Tool].
func (p *PluginTool) GetArgsSchema() string {
	return p.Args
}

// GetBindings implements [tool.Tool].
func (p *PluginTool) GetBindings() tool.Bindings {
	return nil
}

// GetCategories implements [tool.Tool].
func (p *PluginTool) GetCategories() []string {
	return p.Categories
}

// GetDescription implements [tool.Tool].
func (p *PluginTool) GetDescription() string {
	return p.Description
}

// GetName implements [tool.Tool].
func (p *PluginTool) GetName() string {
	return p.Name
}

// WithBindings implements [tool.Tool].
func (p *PluginTool) WithBindings(bindings tool.Bindings) tool.Tool {
	return p
}

func NewPluginTool(name string, args string, description string, categories []string, pluginID string) *PluginTool {
	return &PluginTool{
		Name:        name,
		Args:        args,
		Description: description,
		Categories:  categories,
		PluginID:    pluginID,
	}
}
