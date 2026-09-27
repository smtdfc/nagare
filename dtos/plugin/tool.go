package plugin

import "github.com/smtdfc/nagare/dtos/websocket"

const (
	RegisterPluginToolEvent        websocket.Event = "plugin:register_plugin_tool"
	RegisterPluginToolSuccessEvent websocket.Event = "plugin:register_plugin_tool:success"
	RegisterPluginToolFailedEvent  websocket.Event = "plugin:register_plugin_tool:failed"

	RegisterToolCategoriesEvent        websocket.Event = "plugin:register_tool_categories"
	RegisterToolCategoriesSuccessEvent websocket.Event = "plugin:register_tool_categories:success"
	RegisterToolCategoriesFailedEvent  websocket.Event = "plugin:register_tool_categories:failed"
)

type PluginTool struct {
	Name        string   `json:"name"`
	Args        string   `json:"args"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
}

type RegisterPluginToolEventPayload struct {
	Tool PluginTool `json:"tool"`
}

type RegisterPluginToolSuccessEventPayload struct {
}

type RegisterPluginToolFailedEventPayload struct {
	Cause string `json:"cause"`
}

type RegisterToolCategoriesEventPayload struct {
	Categories []string `json:"categories"`
}

type RegisterToolCategoriesSuccessEventPayload struct {
}

type RegisterToolCategoriesFailedEventPayload struct {
	Cause string `json:"cause"`
}
