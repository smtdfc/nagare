package security

import (
	"encoding/json"
	"testing"
)

func TestAuthPayloadSerialization(t *testing.T) {
	original := AuthPayload{
		ID:         "auth-001",
		TargetType: "plugin",
		Name:       "Calendar Plugin",
		Scopes:     []string{"calendar:read", "calendar:write"},
	}

	// Verify JSON marshaling
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal AuthPayload: %v", err)
	}

	// Verify JSON unmarshaling
	var restored AuthPayload
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal AuthPayload: %v", err)
	}

	if restored.ID != original.ID || restored.TargetType != original.TargetType || restored.Name != original.Name {
		t.Fatalf("payload mismatch after round-trip: %+v", restored)
	}
	if len(restored.Scopes) != len(original.Scopes) {
		t.Fatalf("scope count mismatch: got %d, expected %d", len(restored.Scopes), len(original.Scopes))
	}
}
