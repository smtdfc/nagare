package model_provider

import (
	"github.com/google/uuid"
	core_credential "github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/dtos/rest"
)

func uuidToString(id uuid.UUID) string {
	return id.String()
}

func compatibleToString(c provider.ModelProviderCompatible) string {
	return c.ToString()
}

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
// goverter:extend uuidToString compatibleToString
type IModelProviderDtoMapper interface {
	ToModelProviderDTO(domain *provider.ModelProviderConfig) *rest.ModelProvider
	ToModelProviderDTOs(domains []*provider.ModelProviderConfig) []*rest.ModelProvider
	ToCredentialDTO(credential *core_credential.Credential) *rest.Credential
}
