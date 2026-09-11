package plugin

import (
	"github.com/smtdfc/nagare/shared/dtos/rest"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

const (
	HandshakeEvent                      websocket.Event = "plugin_handshake"
	HandshakeSuccessEvent               websocket.Event = "plugin_handshake_success"
	HandshakeFailedEvent                websocket.Event = "plugin_handshake_failed"
	CreatePluginChatSessionEvent        websocket.Event = "create_plugin_chat_session"
	CreatePluginChatSessionSuccessEvent websocket.Event = "create_plugin_chat_session_success"
	CreatePluginChatSessionErrorEvent   websocket.Event = "create_plugin_chat_session_error"
)

type HandshakeEventPayload struct {
	ID          string `json:"id"`
	PluginID    string `json:"PluginID"`
	ConnectCode string `json:"connectCode"`
}

type HandshakeSuccessEventPayload struct {
	ID string `json:"id"`
}

type HandshakeFailedEventPayload struct {
	ID    string `json:"id"`
	Cause string `json:"cause"`
}

type CreatePluginChatSessionEventPayload struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CreatePluginChatSessionSuccessEventPayload struct {
	ID      string        `json:"id"`
	Session *rest.Session `json:"session"`
}

type CreatePluginChatSessionErrorEventPayload struct {
	ID    string `json:"id"`
	Cause string `json:"cause"`
}
