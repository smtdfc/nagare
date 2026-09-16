package guards

import (
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/custom_errors"
	"github.com/smtdfc/nagare/shared/security"
)

type AuthGuard struct {
	appConfig *config.Config
}

func (a *AuthGuard) VerifyUserFromToken(tokenString string) (*security.AuthPayload, error) {
	auth, err := security.VerifyRSAToken[security.AuthPayload](tokenString, []byte(a.appConfig.PublicKey))
	if err != nil {
		return nil, custom_errors.ErrUnauthorized
	}

	return auth, nil
}

// @Injectable
func NewAuthGuard(
	appConfig *config.Config,
) *AuthGuard {
	return &AuthGuard{
		appConfig: appConfig,
	}
}
