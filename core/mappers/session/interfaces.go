package session

import (
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/session"
)

func parseOwnerType(s string) session.OwnerType {
	return session.GetOwnerType(s)
}

func formatOwnerType(o session.OwnerType) string {
	return string(o)
}

func parseCompatible(s string) provider.ModelProviderCompatible {
	return provider.GetCompatibleFromString(s)
}

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
type ISessionMapper interface {
	// goverter:map OwnerType | parseOwnerType
	ToDomain(entity *entities.Session) *session.Info

	// goverter:map OwnerType | formatOwnerType
	ToEntity(domain *session.Info) *entities.Session

	ToDomains(entities []*entities.Session) []*session.Info

	// goverter:map Compatible | parseCompatible
	ToProviderInfo(entity *entities.ModelProvider) *provider.ModelProviderInfo
}
