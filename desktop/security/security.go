package security

import (
	"github.com/zalando/go-keyring"
)

func GetToken() (string, error) {
	service := "nagare.desktop.keyring"
	token, err := keyring.Get(service, "current_token")
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", err
	}

	return token, nil
}

func SaveToken(token string) error {
	service := "nagare.desktop.keyring"
	return keyring.Set(service, "current_token", token)
}

func ClearToken() error {
	service := "nagare.desktop.keyring"
	return keyring.Delete(service, "current_token")
}
