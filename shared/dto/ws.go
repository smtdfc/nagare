package dto

import "github.com/smtdfc/nagare/shared/messages"

type WsEvent string

const (
	WS_CREATE_SESSION         WsEvent = "WS_CREATE_SESSION"
	WS_CREATE_SESSION_SUCCESS WsEvent = "WS_CREATE_SESSION_SUCCESS"
	WS_CREATE_SESSION_FAILED  WsEvent = "WS_CREATE_SESSION_FAILED"
	WS_INVOKE_AGENT           WsEvent = "WS_INVOKE_AGENT"
	WS_INVOKE_AGENT_FAILED    WsEvent = "WS_INVOKE_AGENT_FAILED"
	WS_AGENT_RESPONSE         WsEvent = "WS_AGENT_RESPONSE"
)

type WsMessage[T any] struct {
	Event   WsEvent `json:"event" mapstructure:"event"`
	Payload T       `json:"payload" mapstructure:"payload"`
}

// Event: WS_CREATE_SESSION_SUCCESS
type CreateSessionSuccess struct {
	ID string `json:"id" mapstructure:"id"`
}

// Event: WS_CREATE_SESSION_FAILED
type CreateSessionFailed struct {
	Cause string `json:"cause" mapstructure:"cause"`
}

// Event: WS_INVOKE_AGENT
type InvokeAgent struct {
	ID        string `json:"id" mapstructure:"id"`
	SessionID string `json:"session_id" mapstructure:"session_id"`
	Text      string `json:"text" mapstructure:"text"`
}

// Event: WS_INVOKE_AGENT_FAILED
type InvokeAgentFailed struct {
	ID    string `json:"id" mapstructure:"id"`
	Cause string `json:"cause" mapstructure:"cause"`
}

// Event: WS_AGENT_RESPONSE
type AgentResponse struct {
	ID        string           `json:"id" mapstructure:"id"`
	SessionID string           `json:"session_id" mapstructure:"session_id"`
	Message   messages.Message `json:"message" mapstructure:"message"`
}
