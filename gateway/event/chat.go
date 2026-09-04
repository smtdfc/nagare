package event

type ChatEventType string

var (
	CHAT_TOPIC                    = "chat"
	CHAT_SEND_EVENT ChatEventType = "CHAT_SEND_EVENT"
)

type ChatEventPayload interface {
	GetEventType() ChatEventType
}

type ChatSendMessageEvent struct {
	SessionID string
	Text      string
}

func (c *ChatSendMessageEvent) GetEventType() ChatEventType {
	return CHAT_SEND_EVENT
}
