package chat

import (
	"context"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/shared/event_bus"
	"github.com/smtdfc/nagare/shared/message"
)

type EventType string

var (
	SendEvent  EventType = "CHAT_SEND_EVENT"
	ChunkEvent EventType = "CHAT_CHUNK_EVENT"
)

type EventPayload interface {
	GetEventType() EventType
}

type TaskInfo struct {
	ID                    string `json:"id"`
	SendResultIntoSession string `json:"sendResultIntoSession"`
}

type SendMessageEventPayload struct {
	RequestID  string
	SessionID  string
	Text       string
	SenderType SenderType
	SenderID   string
	Task       *TaskInfo
}

func (c *SendMessageEventPayload) GetEventType() EventType {
	return SendEvent
}

type ChatChunkEventPayload struct {
	RequestID string
	SessionID string
	Chunk     message.Message
}

func (c *ChatChunkEventPayload) GetEventType() EventType {
	return ChunkEvent
}

type ChatEventBus struct {
	*event_bus.BaseEventBus[EventPayload]
	logger *logger.BaseLogger
}

func (e *ChatEventBus) Subscribe(eventName EventType) (<-chan EventPayload, func()) {
	return e.BaseEventBus.Subscribe(string(eventName))
}

func (e *ChatEventBus) Publish(ctx context.Context, eventName EventType, message EventPayload) {
	e.BaseEventBus.Publish(ctx, string(eventName), message)
}

// @Injectable
func NewEventBus(
	logger *logger.BaseLogger,
) *ChatEventBus {
	bus := event_bus.NewBaseEventBus[EventPayload]()
	return &ChatEventBus{
		BaseEventBus: bus,
		logger:       logger.With("module", "chat-event-bus"),
	}
}
