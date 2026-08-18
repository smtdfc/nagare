package dto

const (
	WS_PLUGIN_REGISTER         WsEvent = "WS_PLUGIN_REGISTER"
	WS_PLUGIN_REGISTER_SUCCESS WsEvent = "WS_PLUGIN_REGISTER_SUCCESS"
	WS_PLUGIN_REGISTER_FAILED  WsEvent = "WS_PLUGIN_REGISTER_FAIL"
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
