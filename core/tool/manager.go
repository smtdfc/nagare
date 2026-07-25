package tool

import (
	"fmt"

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
		return "", fmt.Errorf("Tool not found")
	}

	return tool.Execute(ctx, args)
}

func NewToolManager() *ToolManager {
	return &ToolManager{}
}
