package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

func GenerateRSAToken[T any](payload T, privateKeyBytes []byte, duration time.Duration) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	claims := jwt.MapClaims{
		"user": payload,
		"exp":  time.Now().Add(duration).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func VerifyRSAToken[T any](tokenString string, publicKeyBytes []byte) (*T, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userPayload, ok := claims["user"]
	if !ok {
		return nil, errors.New("user payload not found in token")
	}

	jsonBytes, err := helpers.MarshalJson(userPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payload, err := helpers.UnmarshalJson[T](jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal into target type: %w", err)
	}

	return payload, nil
}
