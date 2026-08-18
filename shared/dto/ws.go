package dto

type WsEvent string

type WsMessage[T any] struct {
	Event   WsEvent `json:"event" mapstructure:"event"`
	Payload T       `json:"payload" mapstructure:"payload"`
}
