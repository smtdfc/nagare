package security

import (
	"encoding/json"

	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
)

func GenerateAuthToken(id string) (string, error) {
	authPayload := dto.AuthPayload{
		ID:   id,
		Role: "user",
	}

	json, err := json.Marshal(&authPayload)
	if err != nil {
		return "", err
	}

	jsonString := string(json)
	keypair, err := GetKeyPair()
	if err != nil {
		return "", err
	}

	return helpers.GenerateToken(keypair.PrivateKeyPEM, jsonString)
}
