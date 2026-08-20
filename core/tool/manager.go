package tool

import (
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
)

type ToolManager struct {
}

func (m *ToolManager) GetListTool() domains.ListTool {
	list := make(domains.ListTool, 0, len(ToolRegistry))
	for _, t := range ToolRegistry {
		list = append(list, t)
	}
	return list
}

func (m *ToolManager) CallTool(ctx domains.Context, name string, args string) (string, error) {
	tool, ok := ToolRegistry[name]
	if !ok {
		return "", custom_errors.NewToolError("Tool resolution failed: The requested tool is either unregistered, deprecated, or unavailable.", name)
	}

	return tool.Execute(ctx, args)
}

// @Injectable
func NewToolManager() *ToolManager {
	return &ToolManager{}
}
