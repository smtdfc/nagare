package event_bus

type EventType string

type EventPayload interface {
	GetEventType() EventType
}
