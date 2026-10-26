package credential

import (
	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
type ICredentialMapper interface {
	ToDomain(entity *entities.Credential) *credential.Credential
	ToEntity(domain *credential.Credential) *entities.Credential
}
