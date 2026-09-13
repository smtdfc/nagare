package websocket

type AuthData struct {
	TargetType string   `json:"targetType"`
	TargetID   string   `json:"targetID"`
	Scopes     []string `json:"scopes"`
}
