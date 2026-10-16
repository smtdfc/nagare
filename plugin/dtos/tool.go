package dtos

const (
	RegisterPluginToolEvent        Event = "plugin:register_plugin_tool"
	RegisterPluginToolSuccessEvent Event = "plugin:register_plugin_tool:success"
	RegisterPluginToolFailedEvent  Event = "plugin:register_plugin_tool:failed"

	RegisterToolCategoriesEvent        Event = "plugin:register_tool_categories"
	RegisterToolCategoriesSuccessEvent Event = "plugin:register_tool_categories:success"
	RegisterToolCategoriesFailedEvent  Event = "plugin:register_tool_categories:failed"
	PluginToolCallEvent                Event = "plugin:tool_call"
	PluginToolCallResultEvent          Event = "plugin:tool_call_result"
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
	Categories map[string]string `json:"categories"`
}

type RegisterToolCategoriesSuccessEventPayload struct {
}

type RegisterToolCategoriesFailedEventPayload struct {
	Cause string `json:"cause"`
}

type PluginToolCallEventPayload struct {
	Name string `json:"name"`
	Args string `json:"args"`
}

type PluginToolCallResultEventPayload struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}
