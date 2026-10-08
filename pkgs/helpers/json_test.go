package helpers

import (
	"testing"
)

type samplePayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestMarshalAndUnmarshalJson(t *testing.T) {
	original := samplePayload{ID: 42, Name: "Nagare"}

	// Test successful marshaling
	raw, err := MarshalJson(original)
	if err != nil {
		t.Fatalf("unexpected error marshaling payload: %v", err)
	}

	// Test successful unmarshaling
	unmarshaled, err := UnmarshalJson[samplePayload](raw)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling payload: %v", err)
	}
	if unmarshaled == nil || unmarshaled.ID != 42 || unmarshaled.Name != "Nagare" {
		t.Fatalf("unexpected unmarshaled content: %+v", unmarshaled)
	}

	// Test unmarshaling invalid JSON string
	invalidRaw := "{invalid_json"
	invalidResult, err := UnmarshalJson[samplePayload](invalidRaw)
	if err == nil {
		t.Fatalf("expected error for malformed json, got nil")
	}
	if invalidResult != nil {
		t.Fatalf("expected nil result pointer on error, got %+v", invalidResult)
	}
}
