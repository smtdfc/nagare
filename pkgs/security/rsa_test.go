package security

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestGenerateRSAKeys(t *testing.T) {
	// Test key generation with 2048 bits
	pair, err := GenerateRSAKeys(2048)
	if err != nil {
		t.Fatalf("unexpected error generating RSA keys: %v", err)
	}

	if pair == nil || len(pair.PublicKey) == 0 || len(pair.PrivateKey) == 0 {
		t.Fatalf("expected non-empty key pair, got %+v", pair)
	}

	// Verify Private Key PEM decoding and structure
	privBlock, _ := pem.Decode(pair.PrivateKey)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		t.Fatalf("invalid private key PEM block: %v", privBlock)
	}
	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse PKCS1 private key: %v", err)
	}
	if privKey.N.BitLen() < 2048 {
		t.Fatalf("expected key size >= 2048, got %d", privKey.N.BitLen())
	}

	// Verify Public Key PEM decoding and structure
	pubBlock, _ := pem.Decode(pair.PublicKey)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		t.Fatalf("invalid public key PEM block: %v", pubBlock)
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse PKIX public key: %v", err)
	}
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("parsed key is not an *rsa.PublicKey")
	}

	// Verify public key matches private key modulus
	if privKey.PublicKey.N.Cmp(pubKey.N) != 0 {
		t.Fatalf("public key does not match private key modulus")
	}

	// Test with invalid bit size
	_, err = GenerateRSAKeys(0)
	if err == nil {
		t.Fatalf("expected error for 0 bits key generation, got nil")
	}
}

func TestRSAPairStructure(t *testing.T) {
	pub := []byte("test-pub")
	priv := []byte("test-priv")
	p := RSAPair{
		PublicKey:  pub,
		PrivateKey: priv,
	}

	if !bytes.Equal(p.PublicKey, pub) || !bytes.Equal(p.PrivateKey, priv) {
		t.Fatalf("unexpected RSAPair values: %+v", p)
	}
}
