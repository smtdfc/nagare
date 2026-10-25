package metadata

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	ArchiveSignatureAlgorithm = "ed25519"
	archiveSignatureDomain    = "nagare-plugin-archive-v1"
)

var ErrInvalidSignature = errors.New("invalid plugin archive signature")

type archiveSignatureFile struct {
	Algorithm string `json:"algorithm"`
	Signature string `json:"signature"`
}

// LoadPrivateKey reads an Ed25519 private key from a PEM file, or from a raw
// base64/hex encoded 32-byte seed or 64-byte private key.
func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key %s: %w", path, err)
	}

	key, err := ParsePrivateKey(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key %s: %w", path, err)
	}

	return key, nil
}

func ParsePrivateKey(raw []byte) (ed25519.PrivateKey, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, errors.New("private key is empty")
	}

	if block, _ := pem.Decode(trimmed); block != nil {
		parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
		}

		key, ok := parsed.(ed25519.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not an Ed25519 key")
		}
		return key, nil
	}

	decoded, err := decodeKeyMaterial(trimmed)
	if err != nil {
		return nil, err
	}

	switch len(decoded) {
	case ed25519.SeedSize:
		return ed25519.NewKeyFromSeed(decoded), nil
	case ed25519.PrivateKeySize:
		return ed25519.PrivateKey(decoded), nil
	default:
		return nil, fmt.Errorf("invalid Ed25519 private key length: got %d bytes, want %d or %d", len(decoded), ed25519.SeedSize, ed25519.PrivateKeySize)
	}
}

// ParsePublicKey accepts a PEM encoded PKIX public key, or a raw base64/hex
// encoded 32-byte Ed25519 public key.
func ParsePublicKey(raw string) (ed25519.PublicKey, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, errors.New("public key is empty")
	}

	if block, _ := pem.Decode([]byte(trimmed)); block != nil {
		parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKIX public key: %w", err)
		}

		key, ok := parsed.(ed25519.PublicKey)
		if !ok {
			return nil, errors.New("public key is not an Ed25519 key")
		}
		return key, nil
	}

	decoded, err := decodeKeyMaterial([]byte(trimmed))
	if err != nil {
		return nil, err
	}
	if len(decoded) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid Ed25519 public key length: got %d bytes, want %d", len(decoded), ed25519.PublicKeySize)
	}

	return ed25519.PublicKey(decoded), nil
}

func decodeKeyMaterial(raw []byte) ([]byte, error) {
	encoded := string(raw)
	var candidates [][]byte

	if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
		candidates = append(candidates, decoded)
	}
	if decoded, err := hex.DecodeString(encoded); err == nil {
		candidates = append(candidates, decoded)
	}

	for _, candidate := range candidates {
		if len(candidate) == ed25519.SeedSize || len(candidate) == ed25519.PrivateKeySize {
			return candidate, nil
		}
	}

	return nil, errors.New("key is not valid base64 or hex, or has an invalid length")
}

// SignArchive signs the SHA-256 digest of the whole plugin archive. Signing the
// archive instead of the binary inside it means metadata and every packaged
// file are covered by the same signature.
func SignArchive(archivePath string, privateKey ed25519.PrivateKey) ([]byte, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid Ed25519 private key")
	}

	digest, err := archiveDigest(archivePath)
	if err != nil {
		return nil, err
	}

	return ed25519.Sign(privateKey, signatureMessage(digest)), nil
}

func VerifyArchive(archivePath string, signature []byte, publicKey ed25519.PublicKey) error {
	if len(publicKey) != ed25519.PublicKeySize {
		return errors.New("invalid Ed25519 public key")
	}
	if len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("%w: invalid signature length %d", ErrInvalidSignature, len(signature))
	}

	digest, err := archiveDigest(archivePath)
	if err != nil {
		return err
	}

	if !ed25519.Verify(publicKey, signatureMessage(digest), signature) {
		return ErrInvalidSignature
	}

	return nil
}

func WriteSignatureFile(path string, signature []byte) error {
	if len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature length %d", len(signature))
	}

	raw, err := json.MarshalIndent(archiveSignatureFile{
		Algorithm: ArchiveSignatureAlgorithm,
		Signature: base64.StdEncoding.EncodeToString(signature),
	}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, append(raw, '\n'), 0664)
}

func ReadSignatureFile(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseSignatureFile(raw)
}

func ParseSignatureFile(raw []byte) ([]byte, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, errors.New("signature file is empty")
	}

	if trimmed[0] == '{' {
		var envelope archiveSignatureFile
		if err := json.Unmarshal(trimmed, &envelope); err != nil {
			return nil, fmt.Errorf("failed to parse signature file: %w", err)
		}
		if envelope.Algorithm != "" && envelope.Algorithm != ArchiveSignatureAlgorithm {
			return nil, fmt.Errorf("unsupported signature algorithm %q", envelope.Algorithm)
		}

		signature, err := base64.StdEncoding.DecodeString(envelope.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to decode signature: %w", err)
		}
		if len(signature) != ed25519.SignatureSize {
			return nil, fmt.Errorf("invalid signature length %d", len(signature))
		}
		return signature, nil
	}

	if signature, err := base64.StdEncoding.DecodeString(string(trimmed)); err == nil && len(signature) == ed25519.SignatureSize {
		return signature, nil
	}
	if signature, err := hex.DecodeString(string(trimmed)); err == nil && len(signature) == ed25519.SignatureSize {
		return signature, nil
	}

	return nil, errors.New("signature is not valid base64 or hex")
}

func archiveDigest(archivePath string) ([]byte, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin archive %s: %w", archivePath, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, fmt.Errorf("failed to hash plugin archive: %w", err)
	}

	return hasher.Sum(nil), nil
}

func signatureMessage(digest []byte) []byte {
	message := make([]byte, 0, len(archiveSignatureDomain)+1+len(digest))
	message = append(message, archiveSignatureDomain...)
	message = append(message, 0)
	message = append(message, digest...)
	return message
}
