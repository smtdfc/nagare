package shared

import "encoding/json"

type Message struct {
	PluginID string          `json:"plugin_id"`
	Kind     string          `json:"method"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}

type RegisterChatChannelPayload struct {
	ID string `json:"id"`
}

type RegisterChatChannelSuccessPayload struct {
	ID string `json:"id"`
}

type HandleChatMessagePayload struct {
	Channel string `json:"id"`
	Message string `json:"message"`
}
