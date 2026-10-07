package provider

import (
	"github.com/google/uuid"
)

type LLMProviderInfo struct {
	ID         uuid.UUID             `json:"id"`
	Name       string                `json:"name"`
	Compatible LLMProviderCompatible `json:"compatible"`
}
