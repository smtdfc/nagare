package plugin

import (
	"github.com/smtdfc/nagare/dtos/websocket"
)

const (
	HandshakeEvent                 websocket.Event = "plugin_handshake"
	HandshakeSuccessEvent          websocket.Event = "plugin_handshake_success"
	HandshakeFailedEvent           websocket.Event = "plugin_handshake_failed"
	PrepareChatSessionEvent        websocket.Event = "plugin_session_prepare"
	PrepareChatSessionSuccessEvent websocket.Event = "plugin_session_prepare_success"
	PrepareChatSessionFailedEvent  websocket.Event = "plugin_session_prepare_failed"
	SendChatMessageEvent           websocket.Event = "plugin_message"
	SendChatMessageSuccessEvent    websocket.Event = "plugin_message_success"
	SendChatMessageFailedEvent     websocket.Event = "plugin_message_failed"
	ResetChatChannelEvent          websocket.Event = "plugin_session_reset"
	ResetChatChannelSuccessEvent   websocket.Event = "plugin_session_reset_success"
	ResetChatChannelFailedEvent    websocket.Event = "plugin_session_reset_failed"
)

type HandshakeEventPayload struct {
	// Deprecated: use RequestID instead
	ID          string `json:"id"`
	PluginID    string `json:"PluginID"`
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
