package tool

import "github.com/smtdfc/nagare/core/domains"

type ToolManager struct {
	Tools []domains.Tool
}

func NewToolManager() *ToolManager {
	return &ToolManager{
		Tools: make([]domains.Tool, 0),
	}
}
