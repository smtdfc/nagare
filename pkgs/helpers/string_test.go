package helpers

import (
	"testing"
)

func TestContainsAnyWord(t *testing.T) {
	// Test matching single keyword
	target := "Nagare is an intelligent desktop agent framework"
	keywords := []string{"agent", "cloud"}
	if !ContainsAnyWord(target, keywords) {
		t.Fatalf("expected true for matched keyword 'agent'")
	}

	// Test case insensitivity
	if !ContainsAnyWord(target, []string{"INTELLIGENT"}) {
		t.Fatalf("expected true with case-insensitive match")
	}

	// Test multi-word keyword phrases split by whitespace
	if !ContainsAnyWord(target, []string{"desktop ui"}) {
		t.Fatalf("expected true for phrase containing 'desktop'")
	}

	// Test non-matching keywords
	if ContainsAnyWord(target, []string{"database", "kubernetes"}) {
		t.Fatalf("expected false for completely unmatched keywords")
	}

	// Test empty or whitespace-only inputs
	if ContainsAnyWord(target, []string{"", "   "}) {
		t.Fatalf("expected false for empty keyword list")
	}

	if ContainsAnyWord("", []string{"agent"}) {
		t.Fatalf("expected false for empty target string")
	}
}
