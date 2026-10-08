package messages

import (
	"testing"
)

func TestMessageTypeToString(t *testing.T) {
	mt := TextMessageType
	if mt.ToString() != "TEXT_MESSAGE" {
		t.Fatalf("expected string 'TEXT_MESSAGE', got %q", mt.ToString())
	}
}
