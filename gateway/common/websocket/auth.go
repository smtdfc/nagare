package websocket

import (
	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/gateway/common/custom_errors"
)

type AuthData struct {
	TargetType string   `json:"targetType"`
	TargetID   string   `json:"targetID"`
	Scopes     []string `json:"scopes"`
}

func GetAuth(s *melody.Session, target string) (*AuthData, error) {
	value, exist := s.Get("auth")
	if !exist {
		return nil, custom_errors.ErrUnauthorized
	}

	auth, ok := value.(*AuthData)
	if !ok || auth.TargetType != target {
		return nil, custom_errors.ErrUnauthorized
	}

	return auth, nil
}
