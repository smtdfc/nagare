package event_bus

import (
	"context"
	"testing"
	"time"
)

type TestEvent struct {
	Message string
}

func TestBaseEventBus_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	eb := NewBaseEventBus[any]()
	ch, unsub := eb.Subscribe("test_event")
	defer unsub()

	expectedPayload := &TestEvent{Message: "Hello Event Bus"}

	eb.Publish(ctx, "test_event", expectedPayload)

	select {
	case payload := <-ch:
		event, ok := payload.(*TestEvent)
		if !ok {
			t.Fatalf("Expected *TestEvent, got %T", payload)
		}
		if event.Message != expectedPayload.Message {
			t.Errorf("Expected message %q, got %q", expectedPayload.Message, event.Message)
		}
	case <-ctx.Done():
		t.Fatal("Timeout: did not receive the event in time")
	}
}
