package plugin

import "github.com/smtdfc/nagare/dtos/websocket"

const (
	RegisterDynamicToolsEvent        websocket.Event = "register_dynamic_tools"
	RegisterDynamicToolsSuccessEvent websocket.Event = "register_dynamic_tools_success"
	RegisterDynamicToolsErrorEvent   websocket.Event = "register_dynamic_tools_error"
)

type DynamicTool struct {
	Name        string `json:"name,omitempty"`
	ArgsSchema  string `json:"argsSchema,omitempty"`
	Description string `json:"description,omitempty"`
}

type RegisterDynamicToolsPayload struct {
	Tools []DynamicTool `json:"tools"`
}

type RegisterDynamicToolsSuccessPayload struct{}

type RegisterDynamicToolsErrorPayload struct {
	Cause string `json:"cause"`
}
