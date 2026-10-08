package security

import (
	"testing"
	"time"
)

func TestGenerateAndVerifyRSAToken(t *testing.T) {
	// Generate key pair for testing
	pair, err := GenerateRSAKeys(2048)
	if err != nil {
		t.Fatalf("failed to generate RSA keys for token test: %v", err)
	}

	payload := AuthPayload{
		ID:         "user-123",
		TargetType: "user",
		Name:       "Test User",
		Scopes:     []string{"read", "write"},
	}

	// Test token generation
	tokenStr, err := GenerateRSAToken(payload, pair.PrivateKey, 2*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate RSA token: %v", err)
	}
	if tokenStr == "" {
		t.Fatalf("expected non-empty token string")
	}

	// Test successful token verification
	verifiedPayload, err := VerifyRSAToken[AuthPayload](tokenStr, pair.PublicKey)
	if err != nil {
		t.Fatalf("failed to verify valid RSA token: %v", err)
	}
	if verifiedPayload == nil || verifiedPayload.ID != payload.ID || verifiedPayload.Name != payload.Name {
		t.Fatalf("verified payload mismatch: %+v", verifiedPayload)
	}
	if len(verifiedPayload.Scopes) != 2 || verifiedPayload.Scopes[0] != "read" {
		t.Fatalf("scopes mismatch: %+v", verifiedPayload.Scopes)
	}

	// Test verification with invalid/wrong public key
	otherPair, err := GenerateRSAKeys(2048)
	if err != nil {
		t.Fatalf("failed to generate second RSA key pair: %v", err)
	}
	badKeyResult, err := VerifyRSAToken[AuthPayload](tokenStr, otherPair.PublicKey)
	if err == nil {
		t.Fatalf("expected error when verifying token with mismatched key, got nil")
	}
	if badKeyResult != nil {
		t.Fatalf("expected nil result on verification error, got %+v", badKeyResult)
	}

	// Test verification with malformed token
	malformedResult, err := VerifyRSAToken[AuthPayload]("invalid.token.string", pair.PublicKey)
	if err == nil {
		t.Fatalf("expected error for malformed token string, got nil")
	}
	if malformedResult != nil {
		t.Fatalf("expected nil result on malformed token error, got %+v", malformedResult)
	}

	// Test expired token verification
	expiredToken, err := GenerateRSAToken(payload, pair.PrivateKey, -1*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}
	expiredResult, err := VerifyRSAToken[AuthPayload](expiredToken, pair.PublicKey)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
	if expiredResult != nil {
		t.Fatalf("expected nil result on expired token, got %+v", expiredResult)
	}

	// Test token generation with invalid private key bytes
	_, err = GenerateRSAToken(payload, []byte("invalid-pem"), time.Hour)
	if err == nil {
		t.Fatalf("expected error generating token with invalid private key, got nil")
	}

	// Test token verification with invalid public key bytes
	_, err = VerifyRSAToken[AuthPayload](tokenStr, []byte("invalid-pem"))
	if err == nil {
		t.Fatalf("expected error verifying token with invalid public key, got nil")
	}
}
