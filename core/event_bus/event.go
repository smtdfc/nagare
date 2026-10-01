package event_bus

import (
	"context"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/pkgs/event_bus"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type SenderType string

const (
	User   SenderType = "user"
	Plugin SenderType = "plugin"
	System SenderType = "system"
)

type EventType string

var (
	SendEvent                 EventType = "event:send"
	ChunkEvent                EventType = "event:chunk"
	RefreshTaskEvent          EventType = "event:refresh_task"
	PluginToolCallEvent       EventType = "event:plugin_tool_call"
	PluginToolCallResultEvent EventType = "event:plugin_tool_call_result"
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
	InvokeID   string
	SessionID  string
	ChannelID  string
	Messages   messages.ListMessage
	SenderType SenderType
	SenderID   string
	Task       *TaskInfo
}

func (c *SendMessageEventPayload) GetEventType() EventType {
	return SendEvent
}

type ChatChunkEventPayload struct {
	RequestID        string
	SessionID        string
	ChannelID        string
	Chunk            messages.Message
	SenderType       SenderType
	SenderID         string
	SessionOwnerID   string
	SessionOwnerType string
}

func (c *ChatChunkEventPayload) GetEventType() EventType {
	return ChunkEvent
}

type RefreshTaskEventPayload struct {
}

func (c *RefreshTaskEventPayload) GetEventType() EventType {
	return RefreshTaskEvent
}

type PluginToolCallEventPayload struct {
	RequestID string
	PluginID  string
	Name      string
	Args      string
}

func (c *PluginToolCallEventPayload) GetEventType() EventType {
	return PluginToolCallEvent
}

type PluginToolCallResultEventPayload struct {
	RequestID string
	Result    string
	Error     string
}

func (c *PluginToolCallResultEventPayload) GetEventType() EventType {
	return PluginToolCallResultEvent
}

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
