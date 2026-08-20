package security

import (
	"encoding/json"

	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
)

func GenerateAuthToken(id string) (string, error) {
	authPayload := dto.JwtAuthPayload{
		ID:     id,
		Role:   "user",
		Kind:   "user",
		Scopes: []string{"basic"},
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
