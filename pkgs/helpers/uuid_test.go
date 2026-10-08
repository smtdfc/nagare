package helpers

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUUID(t *testing.T) {
	// Test format validity
	idStr := GenerateUUID()
	parsed, err := uuid.Parse(idStr)
	if err != nil {
		t.Fatalf("expected valid UUID string, got %q, err=%v", idStr, err)
	}

	// Test uniqueness across sequential calls
	idStr2 := GenerateUUID()
	if idStr == idStr2 {
		t.Fatalf("expected unique UUIDs, got identical strings %q", idStr)
	}
	if parsed.Version() != 4 {
		t.Fatalf("expected UUID version 4, got version %d", parsed.Version())
	}
}
