package chat

import "github.com/smtdfc/nagare/shared/event_bus"

type EventType string

var (
	Topic               = "chat"
	SendEvent EventType = "CHAT_SEND_EVENT"
)

type EventPayload interface {
	GetEventType() EventType
}

type SendMessageEvent struct {
	RequestID string
	SessionID string
	Text      string
}

func (c *SendMessageEvent) GetEventType() EventType {
	return SendEvent
}

type EventBus struct {
	*event_bus.BaseEventBus[EventPayload]
}

// @Injectable
func NewEventBus() *EventBus {
	bus := event_bus.NewBaseEventBus[EventPayload]()
	return &EventBus{
		BaseEventBus: bus,
	}
}
