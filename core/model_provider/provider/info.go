package provider

import (
	"github.com/google/uuid"
)

type ModelProviderInfo struct {
	ID         uuid.UUID               `json:"id"`
	Name       string                  `json:"name"`
	Compatible ModelProviderCompatible `json:"compatible"`
}
