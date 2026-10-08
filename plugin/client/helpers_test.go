package client

import (
	"context"
	"testing"
	"time"

	"github.com/smtdfc/nagare/dtos/websocket"
)

type sampleData struct {
	Message string `json:"message"`
}

func TestGetData(t *testing.T) {
	// Test nil payload
	_, err := GetData[sampleData](nil)
	if err == nil {
		t.Fatalf("expected error for nil payload, got nil")
	}

	// Test valid payload
	payload := &websocket.Payload[any]{
		Event: "test.event",
		Data: map[string]interface{}{
			"message": "hello",
		},
	}

	res, err := GetData[sampleData](payload)
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}
	if res == nil || res.Message != "hello" {
		t.Fatalf("unexpected parsed result: %+v", res)
	}
}

func TestSendAndWaitTimeout(t *testing.T) {
	p := NewPlugin()

	// Connector is not connected, so sendAndWait should return an error
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := sendAndWait[sampleData, sampleData](
		p,
		ctx,
		"send.event",
		map[string]string{"foo": "bar"},
		"fail.event",
		"success.event",
	)
	if err == nil {
		t.Fatalf("expected error when sending without active connection, got nil")
	}
}
