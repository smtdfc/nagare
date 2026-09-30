package security

import (
	"errors"

	"github.com/zalando/go-keyring"
)

func SaveKey(serviceName string, key string, value string) error {
	err := keyring.Set(serviceName, key, value)
	if err != nil {
		return err
	}

	return nil
}

func GetKey(serviceName string, key string) (string, error) {
	value, err := keyring.Get(serviceName, key)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", nil
		}

		return value, err
	}

	return value, nil
}
