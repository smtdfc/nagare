package event_bus

import "github.com/smtdfc/nagare/pkgs/messages"

var (
	SendChatMessageEvent  EventType = "event:chat:send"
	CancelChatEvent       EventType = "event:chat:cancel"
	ChatMessageChunkEvent EventType = "event:chat:chunk"
)

type SendChatMessageEventPayload struct {
	RequestID  string
	InvokeID   string
	SessionID  string
	ChannelID  string
	Messages   messages.ListMessage
	SenderType SenderType
	SenderID   string
	Task       *TaskInfo
}

func (c *SendChatMessageEventPayload) GetEventType() EventType {
	return SendChatMessageEvent
}

type ChatMessageChunkEventPayload struct {
	RequestID        string
	SessionID        string
	ChannelID        string
	Chunk            messages.Message
	SenderType       SenderType
	SenderID         string
	SessionOwnerID   string
	SessionOwnerType string
}

func (c *ChatMessageChunkEventPayload) GetEventType() EventType {
	return ChatMessageChunkEvent
}

type CancelChatEventPayload struct {
	RequestID string
	InvokeID  string
}

func (e *CancelChatEventPayload) GetEventType() EventType {
	return CancelChatEvent
}
