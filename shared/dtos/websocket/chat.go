package websocket

const (
	CHAT_RECEIVED_MESSAGE_EVENT       WebsocketEvent = "CHAT_RECEIVED_MESSAGE_EVENT"
	CHAT_LISTEN_MESSAGE_EVENT         WebsocketEvent = "CHAT_LISTEN_MESSAGE_EVENT"
	CHAT_LISTEN_MESSAGE_SUCCESS_EVENT WebsocketEvent = "CHAT_LISTEN_MESSAGE_SUCCESS_EVENT"
	CHAT_LISTEN_MESSAGE_FAIL_EVENT    WebsocketEvent = "CHAT_LISTEN_MESSAGE_FAIL_EVENT"
)

type ChatReceivedMessageEvent struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type ChatListenMessageEvent struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
}

type ChatListenMessageSuccessEvent struct {
	ID string `json:"id"`
}

type ChatListenMessageFailEvent struct {
	ID    string `json:"id"`
	Cause string `json:"cause"`
}
