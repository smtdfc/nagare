package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/pkgs/security"
)

const SERVICE_NAME = "nagare.security.keyring"

// DefaultTokenTTL is how long a token minted by the CLI stays valid.
const DefaultTokenTTL = 60 * time.Hour

// GenerateUserToken mints a user token signed with the keyring private key.
// It creates a new RSA key pair when none exists yet.
func GenerateUserToken(ttl time.Duration) (string, error) {
	_, privateKey, err := GetRSAKey()
	if err != nil {
		return "", err
	}

	payload := security.AuthPayload{
		ID:         uuid.Nil.String(),
		TargetType: "user",
		Name:       "User",
		Scopes:     []string{},
	}

	return security.GenerateRSAToken(payload, []byte(privateKey), ttl)
}

func GetRSAKey() (string, string, error) {
	privateKey, err := security.GetKey(SERVICE_NAME, "pk")
	if err != nil {
		return "", "", err
	}

	publicKey, err := security.GetKey(SERVICE_NAME, "pubk")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error retrieving public key: %v\n", err)
		return "", "", err
	}

	if privateKey == "" || publicKey == "" {
		rsa, err := security.GenerateRSAKeys(4096)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error generating RSA keys: %v\n", err)
			return "", "", err
		}

		err = security.SaveKey(SERVICE_NAME, "pk", string(rsa.PrivateKey))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error saving private key: %v\n", err)
			return "", "", err
		}

		err = security.SaveKey(SERVICE_NAME, "pubk", string(rsa.PublicKey))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error saving public key: %v\n", err)
			return "", "", err
		}

		privateKey = string(rsa.PrivateKey)
		publicKey = string(rsa.PublicKey)
	}

	return publicKey, privateKey, nil
}
