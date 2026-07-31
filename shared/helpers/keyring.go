package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RSAKeys struct {
	PublicKeyPEM  string
	PrivateKeyPEM string
}

func SignData(privateKeyPEM string, data string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	var privKey *rsa.PrivateKey
	var err error

	if block.Type == "RSA PRIVATE KEY" {
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.New("unsupported private key format: " + block.Type)
		}
		var ok bool
		privKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("key is not an RSA private key")
		}
	}

	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func GenerateToken(privateKeyPEM string, payload string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("invalid private key PEM")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
		var ok bool
		privKey, ok = parsed.(*rsa.PrivateKey)
		if !ok {
			return "", errors.New("not an RSA private key")
		}
	}

	claims := jwt.MapClaims{
		"data": payload,
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(publicKeyPEM string, tokenString string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("invalid public key PEM")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		pubKeyInterface, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return "", errors.New("failed to parse public key")
		}
	}

	pubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload, ok := claims["data"].(string)
		if !ok {
			return "", errors.New("payload is not a valid string")
		}
		return payload, nil
	}

	return "", errors.New("invalid token")
}

func GenerateRSAKeyPair(bits int) (*RSAKeys, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return &RSAKeys{
		PublicKeyPEM:  string(publicKeyPEM),
		PrivateKeyPEM: string(privateKeyPEM),
	}, nil
}
