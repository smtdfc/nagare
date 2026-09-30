package ipc

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestIPCWriteAndRead(t *testing.T) {
	tempDir := t.TempDir()
	sockPath := filepath.Join(tempDir, "test.sock")

	listener, err := Listen(sockPath)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	testMsg := []byte("hello from nagare ipc!")

	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("failed to accept: %v", err)
			return
		}
		defer conn.Close()

		msg, err := ReadMessage(conn)
		if err != nil {
			t.Errorf("failed to read message: %v", err)
			return
		}

		if string(msg) != string(testMsg) {
			t.Errorf("expected %q, got %q", string(testMsg), string(msg))
			return
		}

		resp := []byte("pong response")
		if err := WriteMessage(conn, resp); err != nil {
			t.Errorf("failed to write response: %v", err)
			return
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	clientConn, err := Dial(ctx, sockPath)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer clientConn.Close()

	if err := WriteMessage(clientConn, testMsg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	resp, err := ReadMessage(clientConn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if string(resp) != "pong response" {
		t.Fatalf("expected 'pong response', got %q", string(resp))
	}

	<-done
}
