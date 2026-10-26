package credential

import (
	"github.com/google/uuid"
	core_credential "github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/dtos/rest"
)

func uuidToString(id uuid.UUID) string {
	return id.String()
}

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
// goverter:extend uuidToString
type ICredentialDtoMapper interface {
	ToCredentialDTO(value *core_credential.Credential) *rest.Credential
	ToCredentialDTOs(values []*core_credential.Credential) []*rest.Credential
}
