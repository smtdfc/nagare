package websocket

const (
	ReceivedChatMessageEvent         Event = "CHAT_RECEIVED_MESSAGE_EVENT"
	RegisterChatListenerEvent        Event = "CHAT_LISTEN_MESSAGE_EVENT"
	RegisterChatListenerSuccessEvent Event = "CHAT_LISTEN_MESSAGE_SUCCESS_EVENT"
	RegisterChatListenerFailEvent    Event = "CHAT_LISTEN_MESSAGE_FAIL_EVENT"
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
