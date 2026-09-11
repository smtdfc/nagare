package websocket

const (
	ReceivedChatMessageEvent         Event = "ReceivedChatMessageEvent"
	RegisterChatListenerEvent        Event = "RegisterChatListenerEvent"
	RegisterChatListenerSuccessEvent Event = "RegisterChatListenerSuccessEvent"
	RegisterChatListenerFailEvent    Event = "RegisterChatListenerFailEvent"
)

type ReceivedChatMessageEventPayload struct {
	SessionID string `json:"sessionID"`
	Message   string `json:"message"`
}

type RegisterChatMessageListenerEventPayload struct {
	ID        string `json:"id"`
	SessionID string `json:"sessionID"`
}

type RegisterChatListenerSuccessEventEventPayload struct {
	ID string `json:"id"`
}

type RegisterChatListenerFailEventPayload struct {
	ID    string `json:"id"`
	Cause string `json:"cause"`
}
