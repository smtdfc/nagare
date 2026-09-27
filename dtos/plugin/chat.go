package plugin

import (
	"github.com/smtdfc/nagare/dtos/websocket"
)

const (
	HandshakeEvent        websocket.Event = "plugin:handshake"
	HandshakeSuccessEvent websocket.Event = "plugin:handshake:success"
	HandshakeFailedEvent  websocket.Event = "plugin:handshake:failed"

	PrepareChatSessionEvent        websocket.Event = "plugin:prepare_chat_session"
	PrepareChatSessionSuccessEvent websocket.Event = "plugin:prepare_chat_session:success"
	PrepareChatSessionFailedEvent  websocket.Event = "plugin:prepare_chat_session:failed"

	SendChatMessageEvent        websocket.Event = "plugin:send_chat_message"
	SendChatMessageSuccessEvent websocket.Event = "plugin:send_chat_message:success"
	SendChatMessageFailedEvent  websocket.Event = "plugin:send_chat_message:failed"

	ResetChatChannelEvent        websocket.Event = "plugin:reset_chat_channel"
	ResetChatChannelSuccessEvent websocket.Event = "plugin:reset_chat_channel:success"
	ResetChatChannelFailedEvent  websocket.Event = "plugin:reset_chat_channel:failed"
)

type HandshakeEventPayload struct {
	// Deprecated: use RequestID instead
	ID          string `json:"id"`
	PackageName string `json:"packageName"`
	ConnectCode string `json:"connectCode"`
}

type HandshakeSuccessEventPayload struct {
	// Deprecated: use RequestID instead
	ID string `json:"id"`
}

type HandshakeFailedEventPayload struct {
	// Deprecated: use RequestID instead
	ID    string `json:"id"`
	Cause string `json:"cause"`
}

type PrepareChatSessionEventPayload struct {
	ChannelID string `json:"channelID"`
}

type PrepareChatSessionSuccessEventPayload struct {
	ChannelID string `json:"channelID"`
	SessionID string `json:"sessionID"`
}

type PrepareChatSessionFailedEventPayload struct {
	ChannelID string `json:"channelID"`
	Cause     string `json:"cause"`
}

type SendChatMessageEventPayload struct {
	SessionID string `json:"sessionID"`
	Text      string `json:"text"`
}

type SendChatMessageSuccessEventPayload struct {
	SessionID string `json:"sessionID"`
}

type SendChatMessageFailedEventPayload struct {
	Cause string `json:"cause"`
}

type ResetChatChannelEventPayload struct {
	ChannelID string `json:"channelID"`
}

type ResetChatChannelSuccessEventPayload struct {
	ChannelID string `json:"channelID"`
}

type ResetChatChannelFailedEventPayload struct {
	Cause string `json:"cause"`
}
