package security

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/zalando/go-keyring"
)

func SaveKey(serviceName string, key string, value string) error {
	if runtime.GOOS == "android" {
		return saveTextKey(serviceName, key, value)
	}

	return keyring.Set(serviceName, key, value)
}

func GetKey(serviceName string, key string) (string, error) {
	if runtime.GOOS == "android" {
		return getTextKey(serviceName, key)
	}

	value, err := keyring.Get(serviceName, key)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", nil
	}

	return value, err
}

func saveTextKey(serviceName, key, value string) error {
	path := filepath.Join(paths.CredentialsDir, serviceName+"_"+key)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(value), 0600)
}

func getTextKey(serviceName, key string) (string, error) {
	path := filepath.Join(paths.CredentialsDir, serviceName+"_"+key)

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return string(data), nil
}
