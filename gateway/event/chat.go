package event

type ChatEventType string

var (
	ChatTopic                   = "chat"
	ChatSendEvent ChatEventType = "CHAT_SEND_EVENT"
)

type ChatEventPayload interface {
	GetEventType() ChatEventType
}

type ChatSendMessageEvent struct {
	SessionID string
	Text      string
}

func (c *ChatSendMessageEvent) GetEventType() ChatEventType {
	return ChatSendEvent
}
