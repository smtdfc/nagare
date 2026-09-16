package websocket

const (
	AuthEvent        Event = "AuthEvent"
	AuthSuccessEvent Event = "AuthSuccessEvent"
	AuthFailedEvent  Event = "AuthFailedEvent"
)

type AuthEventPayload struct {
	Token string `json:"token"`
}

type AuthSuccessEventPayload struct{}
type AuthFailedEventPayload struct {
	Cause string `json:"cause"`
}
