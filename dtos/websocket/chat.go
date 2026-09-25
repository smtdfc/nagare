package websocket

const (
	ReceivedChatMessageEvent           Event = "ReceivedChatMessageEvent"
	RegisterChatListenerEvent          Event = "RegisterChatListenerEvent"
	RegisterChatListenerSuccessEvent   Event = "RegisterChatListenerSuccessEvent"
	RegisterChatListenerFailEvent      Event = "RegisterChatListenerFailEvent"
	UnregisterChatListenerEvent        Event = "UnregisterChatListenerEvent"
	UnregisterChatListenerSuccessEvent Event = "UnregisterChatListenerSuccessEvent"
	UnregisterChatListenerFailEvent    Event = "UnregisterChatListenerFailEvent"
)

type ReceivedChatMessageEventPayload struct {
	SessionID string `json:"sessionID"`
	ChannelID string `json:"channelID"`
	Message   string `json:"message"`
}

type RegisterChatMessageListenerEventPayload struct {
	// Deprecated: use RequestID instead
	ID        string `json:"id"`
	SessionID string `json:"sessionID"`
}

type RegisterChatListenerSuccessEventEventPayload struct {
	// Deprecated: use RequestID instead
	ID string `json:"id"`
}

type RegisterChatListenerFailEventPayload struct {
	// Deprecated: use RequestID instead
	ID    string `json:"id"`
	Cause string `json:"cause"`
}
