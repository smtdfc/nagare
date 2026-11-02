package external_module

type Event string

type Payload[T any] struct {
	RequestID string `json:"requestID"`
	Event     Event  `json:"event"`
	Data      T      `json:"data"`
}
