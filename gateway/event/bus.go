package event

import "github.com/smtdfc/nagare/shared/event_bus"

type AppEventBusSystem struct {
	ChatEventBus *event_bus.EventBus[ChatEventPayload]
}

// @Injectable
func NewAppEventBus() *AppEventBusSystem {
	return &AppEventBusSystem{
		ChatEventBus: event_bus.NewEventBus[ChatEventPayload](),
	}
}
