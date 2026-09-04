package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smtdfc/nagare/shared/helpers"
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
	var zero = new(T)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return zero, fmt.Errorf("failed to parse public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil {
		return zero, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return zero, errors.New("invalid token")
	}

	userPayload, ok := claims["user"]
	if !ok {
		return zero, errors.New("user payload not found in token")
	}

	jsonBytes, err := helpers.MarshalJson(userPayload)
	if err != nil {
		return zero, fmt.Errorf("failed to marshal user payload: %w", err)
	}

	payload, err := helpers.UnmarshalJson[T](jsonBytes)
	if err != nil {
		return zero, fmt.Errorf("failed to unmarshal into target type: %w", err)
	}

	return payload, nil
}
