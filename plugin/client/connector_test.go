package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/smtdfc/nagare/dtos/websocket"
)

func TestConnectorSendWithoutConnection(t *testing.T) {
	connector := NewConnector(func(p *websocket.Payload[any]) {})
	err := connector.Send("test.event", "data", "req-1")
	if !errors.Is(err, ErrConnectionNotReady) {
		t.Fatalf("expected ErrConnectionNotReady, got %v", err)
	}

	// Close when already closed should return nil
	if err := connector.Close(); err != nil {
		t.Fatalf("unexpected error closing unconnected connector: %v", err)
	}
}

func TestConnectorWaitCancelled(t *testing.T) {
	connector := NewConnector(nil)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := connector.Wait(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}
