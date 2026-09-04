package websocket

type WebsocketEvent string
type WebsocketPayload[T any] struct {
	Event WebsocketEvent `json:"event"`
	Data  T              `json:"data"`
}
