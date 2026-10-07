package credential

import "github.com/google/uuid"

type Credential struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	ApiKey string    `json:"api_key"`
}
