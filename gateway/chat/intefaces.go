package chat

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/session"
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
type ISessionDtoMapper interface {
	// goverter:map ModelProvider CurrentModelProvider
	ToSessionDTO(s *session.Info) *rest.Session
	ToSessionDTOs(sessions []*session.Info) []*rest.Session

	ToModelProviderInfoDTO(p *provider.ModelProviderInfo) *rest.ModelProviderInfo
}
