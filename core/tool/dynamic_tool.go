package tool

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/plugin"
)

type DynamicTool struct {
	Name        string
	Description string
	Args        string
	PluginID    uuid.UUID
	Plugin      *plugin.Plugin
}

type ToolMetadata struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Args         string `json:"args"`
	IsPluginTool bool   `json:"is_plugin_tool"`
	PluginID     string `json:"plugin_id,omitempty"`
}
