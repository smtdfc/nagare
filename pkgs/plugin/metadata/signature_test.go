package metadata

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, name string, content []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
	return path
}

func TestSignAndVerifyArchive(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	archivePath := writeTestFile(t, "plugin.nagare_plugin", []byte("archive content"))

	signature, err := SignArchive(archivePath, privateKey)
	if err != nil {
		t.Fatalf("failed to sign archive: %v", err)
	}

	if err := VerifyArchive(archivePath, signature, publicKey); err != nil {
		t.Fatalf("expected valid signature, got: %v", err)
	}

	if err := VerifyArchive(archivePath, signature, privateKey.Public().(ed25519.PublicKey)); err != nil {
		t.Fatalf("expected valid signature with derived public key, got: %v", err)
	}
}

func TestVerifyArchiveRejectsTampering(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	archivePath := writeTestFile(t, "plugin.nagare_plugin", []byte("archive content"))

	signature, err := SignArchive(archivePath, privateKey)
	if err != nil {
		t.Fatalf("failed to sign archive: %v", err)
	}

	if err := os.WriteFile(archivePath, []byte("tampered content"), 0644); err != nil {
		t.Fatalf("failed to tamper archive: %v", err)
	}

	if err := VerifyArchive(archivePath, signature, publicKey); !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature for tampered archive, got: %v", err)
	}
}

func TestVerifyArchiveRejectsWrongKey(t *testing.T) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	otherPublicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate second key: %v", err)
	}

	archivePath := writeTestFile(t, "plugin.nagare_plugin", []byte("archive content"))

	signature, err := SignArchive(archivePath, privateKey)
	if err != nil {
		t.Fatalf("failed to sign archive: %v", err)
	}

	if err := VerifyArchive(archivePath, signature, otherPublicKey); !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature for wrong key, got: %v", err)
	}
}

func TestSignatureFileRoundTrip(t *testing.T) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	archivePath := writeTestFile(t, "plugin.nagare_plugin", []byte("archive content"))
	signature, err := SignArchive(archivePath, privateKey)
	if err != nil {
		t.Fatalf("failed to sign archive: %v", err)
	}

	signaturePath := filepath.Join(filepath.Dir(archivePath), "plugin.nagare_plugin.sig")
	if err := WriteSignatureFile(signaturePath, signature); err != nil {
		t.Fatalf("failed to write signature file: %v", err)
	}

	parsed, err := ReadSignatureFile(signaturePath)
	if err != nil {
		t.Fatalf("failed to read signature file: %v", err)
	}
	if string(parsed) != string(signature) {
		t.Fatalf("signature round trip mismatch")
	}

	raw := []byte(base64.StdEncoding.EncodeToString(signature))
	parsedRaw, err := ParseSignatureFile(raw)
	if err != nil {
		t.Fatalf("failed to parse raw base64 signature: %v", err)
	}
	if string(parsedRaw) != string(signature) {
		t.Fatalf("raw base64 signature round trip mismatch")
	}
}

func TestParsePrivateKeyFormats(t *testing.T) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	seed := privateKey.Seed()

	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	cases := []struct {
		name string
		raw  []byte
	}{
		{name: "pem", raw: pemBytes},
		{name: "base64 seed", raw: []byte(base64.StdEncoding.EncodeToString(seed))},
		{name: "hex seed", raw: []byte(strings.ToUpper(hex.EncodeToString(seed)))},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := ParsePrivateKey(tc.raw)
			if err != nil {
				t.Fatalf("failed to parse key: %v", err)
			}
			if !parsed.Equal(privateKey) {
				t.Fatalf("parsed key does not match original")
			}
		})
	}
}

func TestParsePublicKeyFormats(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})

	cases := []struct {
		name string
		raw  string
	}{
		{name: "pem", raw: string(pemBytes)},
		{name: "base64", raw: base64.StdEncoding.EncodeToString(publicKey)},
		{name: "hex", raw: hex.EncodeToString(publicKey)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := ParsePublicKey(tc.raw)
			if err != nil {
				t.Fatalf("failed to parse key: %v", err)
			}
			if !parsed.Equal(publicKey) {
				t.Fatalf("parsed key does not match original")
			}
		})
	}
}

func TestParsePrivateKeyRejectsInvalidInput(t *testing.T) {
	if _, err := ParsePrivateKey(nil); err == nil {
		t.Fatalf("expected error for empty key")
	}
	if _, err := ParsePrivateKey([]byte("not-a-key")); err == nil {
		t.Fatalf("expected error for invalid key")
	}
}
