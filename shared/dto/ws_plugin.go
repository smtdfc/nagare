package dto

const (
	WS_PLUGIN_REGISTER             WsEvent = "WS_PLUGIN_REGISTER"
	WS_PLUGIN_REGISTER_SUCCESS     WsEvent = "WS_PLUGIN_REGISTER_SUCCESS"
	WS_PLUGIN_REGISTER_FAILED      WsEvent = "WS_PLUGIN_REGISTER_FAIL"
	WS_PLUGIN_INVOKE_AGENT         WsEvent = "WS_PLUGIN_INVOKE_AGENT"
	WS_PLUGIN_INVOKE_AGENT_SUCCESS WsEvent = "WS_PLUGIN_INVOKE_AGENT_SUCCESS"
	WS_PLUGIN_INVOKE_AGENT_FAILED  WsEvent = "WS_PLUGIN_INVOKE_AGENT_FAILED"
)

// Event: WS_PLUGIN_REGISTER
type PluginRegisterEvent struct {
	Code string `json:"code" mapstructure:"code"`
}

// Event: WS_PLUGIN_REGISTER_SUCCESS
type PluginRegisterSuccessEvent struct {
}

// Event: WS_PLUGIN_REGISTER_FAILED
type PluginRegisterFailedEvent struct {
	Cause string `json:"cause" mapstructure:"cause"`
}

// Event: WS_PLUGIN_INVOKE_AGENT
type PluginInvokeAgentEvent struct {
	ID        string `json:"id" mapstructure:"id"`
	SessionID string `json:"session_id" mapstructure:"session_id"`
	Text      string `json:"text" mapstructure:"text"`
}

// Event: WS_PLUGIN_INVOKE_AGENT_FAILED
type PluginInvokeAgentFailedEvent struct {
	ID    string `json:"id" mapstructure:"id"`
	Cause string `json:"cause" mapstructure:"cause"`
}

// Event: WS_PLUGIN_INVOKE_AGENT_SUCCESS
type PluginInvokeAgentSuccessEvent struct {
	ID        string `json:"id" mapstructure:"id"`
	SessionID string `json:"session_id" mapstructure:"session_id"`
	Message   any    `json:"message" mapstructure:"message"`
}
