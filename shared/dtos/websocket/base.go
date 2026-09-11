package websocket

type Event string
type Payload[T any] struct {
	Event Event `json:"event"`
	Data  T     `json:"data"`
}
