package helpers

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/shared/security"
)

const SERVICE_NAME = "nagare.security.keyring"

func GetRSAKey() (string, string, error) {
	privateKey, err := security.GetKey(SERVICE_NAME, "pk")
	if err != nil {
		return "", "", err
	}

	publicKey, err := security.GetKey(SERVICE_NAME, "pubk")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving public key: %v\n", err)
		return "", "", err
	}

	if privateKey == "" || publicKey == "" {
		rsa, err := security.GenerateRSAKeys(4096)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating RSA keys: %v\n", err)
			return "", "", err
		}

		err = security.SaveKey(SERVICE_NAME, "pk", string(rsa.PrivateKey))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving private key: %v\n", err)
			return "", "", err
		}

		err = security.SaveKey(SERVICE_NAME, "pubk", string(rsa.PublicKey))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving public key: %v\n", err)
			return "", "", err
		}

		privateKey = string(rsa.PrivateKey)
		publicKey = string(rsa.PublicKey)
	}

	return publicKey, privateKey, nil
}
