package dto

type WsEvent string

const (
	WS_CREATE_SESSION         WsEvent = "WS_CREATE_SESSION"
	WS_CREATE_SESSION_SUCCESS WsEvent = "WS_CREATE_SESSION_SUCCESS"
	WS_CREATE_SESSION_FAILED  WsEvent = "WS_CREATE_SESSION_FAILED"
)

type WsMessage[T any] struct {
	Event   WsEvent `json:"event"`
	Payload T       `json:"payload"`
}

// Event: WS_CREATE_SESSION_SUCCESS
type CreateSessionSuccess struct {
	ID string `json:"id"`
}

// Event: WS_CREATE_SESSION_FAILED
type CreateSessionFailed struct {
	Cause string `json:"cause"`
}
