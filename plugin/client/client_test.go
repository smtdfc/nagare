package client

import (
	"context"
	"testing"
	"time"

	"github.com/smtdfc/nagare/dtos/websocket"
)

func TestPluginClientNewAndHandleEvent(t *testing.T) {
	p := NewPlugin()
	if p == nil {
		t.Fatalf("expected non-nil PluginClient")
	}

	received := false
	p.OnReceivedChatMessage = func(sessionID, channelID, chunk string) {
		if sessionID == "sess-1" && channelID == "chan-1" && chunk == "test chunk" {
			received = true
		}
	}

	// Test dispatching ReceivedChatMessageEvent
	payload := &websocket.Payload[any]{
		Event: websocket.ReceivedChatMessageEvent,
		Data: map[string]interface{}{
			"sessionID": "sess-1",
			"channelID": "chan-1",
			"message":   "test chunk",
		},
	}
	p.handleEvent(payload)

	if !received {
		t.Fatalf("expected OnReceivedChatMessage to be invoked")
	}

	// Test Start with invalid socket path returns error
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	p.ConnectConfig.SocketPath = "/non/existent/path/socket.sock"
	err := p.Start(ctx, func() {})
	if err == nil {
		t.Fatalf("expected start error on invalid socket path, got nil")
	}
}
