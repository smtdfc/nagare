package shared

import "encoding/json"

type Message struct {
	PluginID string          `json:"plugin_id"`
	Kind     string          `json:"method"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}
