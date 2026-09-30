package security

type AuthPayload struct {
	ID         string   `json:"id"`
	TargetType string   `json:"target_type"`
	Name       string   `json:"name"`
	Scopes     []string `json:"scopes"`
}
