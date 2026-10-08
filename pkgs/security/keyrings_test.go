package security

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/pkgs/paths"
)

func TestTextKeyStorage(t *testing.T) {
	tempDir := t.TempDir()
	origCredsDir := paths.CredentialsDir
	paths.CredentialsDir = tempDir
	defer func() {
		paths.CredentialsDir = origCredsDir
	}()

	service := "test-service"
	key := "api_token"
	val := "secret-token-value-123"

	// Test saving key to text file
	if err := saveTextKey(service, key, val); err != nil {
		t.Fatalf("failed to save text key: %v", err)
	}

	expectedPath := filepath.Join(tempDir, service+"_"+key)
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatalf("expected key file to exist at %s: %v", expectedPath, err)
	}

	// Test reading back the saved key
	readVal, err := getTextKey(service, key)
	if err != nil {
		t.Fatalf("failed to get text key: %v", err)
	}
	if readVal != val {
		t.Fatalf("expected retrieved key %q, got %q", val, readVal)
	}

	// Test retrieving non-existent key returns empty string and nil error
	nonExistentVal, err := getTextKey(service, "non_existent_key")
	if err != nil {
		t.Fatalf("unexpected error for non-existent key: %v", err)
	}
	if nonExistentVal != "" {
		t.Fatalf("expected empty string for non-existent key, got %q", nonExistentVal)
	}
}
