package dto

import "github.com/smtdfc/nagare/shared/messages"

const (
	WS_CREATE_SESSION         WsEvent = "WS_CREATE_SESSION"
	WS_CREATE_SESSION_SUCCESS WsEvent = "WS_CREATE_SESSION_SUCCESS"
	WS_CREATE_SESSION_FAILED  WsEvent = "WS_CREATE_SESSION_FAILED"
	WS_INVOKE_AGENT           WsEvent = "WS_INVOKE_AGENT"
	WS_INVOKE_AGENT_FAILED    WsEvent = "WS_INVOKE_AGENT_FAILED"
	WS_AGENT_RESPONSE         WsEvent = "WS_AGENT_RESPONSE"
	WS_AUTH_REQUEST           WsEvent = "WS_AUTH_REQUEST"
	WS_AUTH_SUCCESS           WsEvent = "WS_AUTH_SUCCESS"
	WS_AUTH_FAILED            WsEvent = "WS_AUTH_FAILED"
)

// Event: WS_CREATE_SESSION_SUCCESS
type CreateSessionSuccessEvent struct {
	ID string `json:"id" mapstructure:"id"`
}

// Event: WS_CREATE_SESSION_FAILED
type CreateSessionFailedEvent struct {
	Cause string `json:"cause" mapstructure:"cause"`
}

// Event: WS_INVOKE_AGENT
type InvokeAgentEvent struct {
	ID        string `json:"id" mapstructure:"id"`
	SessionID string `json:"session_id" mapstructure:"session_id"`
	Text      string `json:"text" mapstructure:"text"`
}

// Event: WS_INVOKE_AGENT_FAILED
type InvokeAgentFailedEvent struct {
	ID    string `json:"id" mapstructure:"id"`
	Cause string `json:"cause" mapstructure:"cause"`
}

// Event: WS_AGENT_RESPONSE
type AgentOutputEvent struct {
	ID        string           `json:"id" mapstructure:"id"`
	SessionID string           `json:"session_id" mapstructure:"session_id"`
	Message   messages.Message `json:"message" mapstructure:"message"`
}

// Event: WS_AUTH
type AuthRequestEvent struct {
	Token string `json:"token" mapstructure:"token"`
}

// Event: WS_AUTH_SUCCESS
type WsAuthSuccessEvent struct {
}

// Event: WS_AUTH_FAILED
type WsAuthFailedEvent struct {
	Cause string `json:"cause" mapstructure:"cause"`
}
