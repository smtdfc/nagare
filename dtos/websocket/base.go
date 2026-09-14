package websocket

type Event string
type Payload[T any] struct {
	Event     Event  `json:"event"`
	RequestID string `json:"requestID"`
	Data      T      `json:"data"`
}
