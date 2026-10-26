package model_provider

import (
	"strings"

	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func splitModels(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func joinModels(models []string) string {
	return strings.Join(models, ",")
}

func parseCompatible(s string) provider.ModelProviderCompatible {
	return provider.GetCompatibleFromString(s)
}

func formatCompatible(c provider.ModelProviderCompatible) string {
	return c.ToString()
}

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
type IModelProviderMapper interface {
	// goverter:map Models | splitModels
	// goverter:map Compatible | parseCompatible
	ToDomain(entity *entities.ModelProvider) *provider.ModelProviderConfig

	// goverter:map Models | joinModels
	// goverter:map Compatible | formatCompatible
	ToEntity(domain *provider.ModelProviderConfig) *entities.ModelProvider

	// goverter:map Compatible | parseCompatible
	ToProviderInfo(entity *entities.ModelProvider) *provider.ModelProviderInfo

	ToDomains(entities []*entities.ModelProvider) []*provider.ModelProviderConfig

	ToCredentialDomain(entity *entities.Credential) *credential.Credential
}
