package provider

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/credential"
)

type ModelProviderCompatible string

const (
	OpenAICompatible  ModelProviderCompatible = "OpenAI"
	UnknownCompatible ModelProviderCompatible = "Unknown"
)

func (l ModelProviderCompatible) ToString() string {
	return string(l)
}

func GetCompatibleFromString(s string) ModelProviderCompatible {
	switch s {
	case string(OpenAICompatible):
		return OpenAICompatible
	default:
		return UnknownCompatible
	}
}

type ModelProviderConfig struct {
	ID         uuid.UUID
	Name       string
	Compatible ModelProviderCompatible
	ApiKey     string
	Models     []string
	BaseURL    string

	CredentialID uuid.UUID
	Credential   *credential.Credential
}
