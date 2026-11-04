package provider

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/credential"
)

type ModelProviderCompatible string

const (
	OpenAICompatible  ModelProviderCompatible = "OpenAI"
	GoogleCompatible  ModelProviderCompatible = "Google"
	UnknownCompatible ModelProviderCompatible = "Unknown"
)

func (l ModelProviderCompatible) ToString() string {
	return string(l)
}

func GetCompatibleFromString(s string) ModelProviderCompatible {
	switch s {
	case string(OpenAICompatible):
		return OpenAICompatible
	case string(GoogleCompatible):
		return GoogleCompatible
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
