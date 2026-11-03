package client

import (
	"context"
	"testing"
)

func TestChatClientOperationsWithoutConnection(t *testing.T) {
	p := NewPlugin()
	ctx := context.Background()

	// Test PrepareChatChannel failure without connection
	err := p.PrepareChatChannel(ctx, "chat-123")
	if err == nil {
		t.Fatalf("expected error when preparing chat channel without connection, got nil")
	}

	// Test SendChatMessage failure without connection
	err = p.SendChatMessage(ctx, "chat-123", "hello")
	if err == nil {
		t.Fatalf("expected error when sending chat message without connection, got nil")
	}

	// Test ResetChatChannel failure without connection
	err = p.ResetChatChannel(ctx, "chat-123")
	if err == nil {
		t.Fatalf("expected error when resetting chat channel without connection, got nil")
	}
}
