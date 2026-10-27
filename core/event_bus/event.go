package event_bus

import (
	"context"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/pkgs/event_bus"
)

type CoreEventBus struct {
	*event_bus.BaseEventBus[EventPayload]
	logger *logger.BaseLogger
}

func (e *CoreEventBus) Subscribe(eventName EventType) (<-chan EventPayload, func()) {
	return e.BaseEventBus.Subscribe(string(eventName))
}

func (e *CoreEventBus) Publish(ctx context.Context, eventName EventType, message EventPayload) {
	e.BaseEventBus.Publish(ctx, string(eventName), message)
}

// @Injectable
func NewEventBus(
	logger *logger.BaseLogger,
) *CoreEventBus {
	bus := event_bus.NewBaseEventBus[EventPayload]()
	return &CoreEventBus{
		BaseEventBus: bus,
		logger:       logger.With("module", "chat-event-bus"),
	}
}
