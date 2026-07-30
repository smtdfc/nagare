package security

import (
	"errors"

	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/zalando/go-keyring"
)

func GetKeyPair() (*helpers.RSAKeys, error) {
	service := "nagare.app.keyring"
	privateKey, errPriv := keyring.Get(service, "app_private")
	publicKey, errPub := keyring.Get(service, "app_public")

	if errors.Is(errPriv, keyring.ErrNotFound) || errors.Is(errPub, keyring.ErrNotFound) || privateKey == "" || publicKey == "" {
		keypair, err := helpers.GenerateRSAKeyPair(2048)
		if err != nil {
			return nil, err
		}

		err = keyring.Set(service, "app_public", keypair.PublicKeyPEM)
		if err != nil {
			return nil, err
		}

		err = keyring.Set(service, "app_private", keypair.PrivateKeyPEM)
		if err != nil {
			return nil, err
		}

		return keypair, nil
	}

	if errPriv != nil {
		return nil, errPriv
	}
	if errPub != nil {
		return nil, errPub
	}

	return &helpers.RSAKeys{
		PublicKeyPEM:  publicKey,
		PrivateKeyPEM: privateKey,
	}, nil
}
