package main

import (
	"testing"
)

func TestMessageBuffer(t *testing.T) {
	mb := &MessageBuffer{}

	// Test appending chunks
	mb.Append("Hello ")
	mb.Append("Telegram ")
	mb.Append("World")

	expected := "Hello Telegram World"
	if mb.String() != expected {
		t.Fatalf("expected buffer content %q, got %q", expected, mb.String())
	}

	// Test concurrent access
	tp := &TelegramPlugin{}
	b1 := tp.getOrCreateBuffer("sess-1")
	b2 := tp.getOrCreateBuffer("sess-1")
	if b1 != b2 {
		t.Fatalf("expected identical buffer instance for same sessionID")
	}

	tp.clearBuffer("sess-1")
	b3 := tp.getOrCreateBuffer("sess-1")
	if b1 == b3 && b3.String() != "" {
		t.Fatalf("expected fresh buffer after clearBuffer")
	}
}
