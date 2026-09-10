package websocket

const (
	ReceivedChatMessageEvent         WebsocketEvent = "CHAT_RECEIVED_MESSAGE_EVENT"
	RegisterChatListenerEvent        WebsocketEvent = "CHAT_LISTEN_MESSAGE_EVENT"
	RegisterChatListenerSuccessEvent WebsocketEvent = "CHAT_LISTEN_MESSAGE_SUCCESS_EVENT"
	RegisterChatListenerFailEvent    WebsocketEvent = "CHAT_LISTEN_MESSAGE_FAIL_EVENT"
)

type ReceivedChatMessageEventPayload struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type RegisterChatMessageListenerEventPayload struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
}

type RegisterChatListenerSuccessEventEventPayload struct {
	ID string `json:"id"`
}

type RegisterChatListenerFailEventPayload struct {
	ID    string `json:"id"`
	Cause string `json:"cause"`
}
