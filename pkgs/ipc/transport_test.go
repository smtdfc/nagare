package ipc

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

type sampleIPCMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func TestTransportFraming(t *testing.T) {
	var buf bytes.Buffer
	testMsg := []byte("Hello Nagare IPC Framing")

	// Test writing framed message
	if err := WriteMessage(&buf, testMsg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	// Test reading framed message back
	readBack, err := ReadMessage(&buf)
	if err != nil {
		t.Fatalf("failed to read framed message: %v", err)
	}

	if !bytes.Equal(readBack, testMsg) {
		t.Fatalf("content mismatch: got %q, expected %q", string(readBack), string(testMsg))
	}

	// Test message exceeding MaxMessageSize
	oversizedMsg := make([]byte, MaxMessageSize+1)
	err = WriteMessage(&buf, oversizedMsg)
	if !errors.Is(err, ErrMessageTooLarge) {
		t.Fatalf("expected ErrMessageTooLarge, got %v", err)
	}

	// Test reading EOF / short read
	emptyBuf := bytes.NewReader([]byte{})
	_, err = ReadMessage(emptyBuf)
	if err != io.EOF {
		t.Fatalf("expected EOF on empty reader, got %v", err)
	}
}

func TestSendJSON(t *testing.T) {
	var buf bytes.Buffer
	msg := sampleIPCMessage{
		Event: "handshake",
		Data:  "ok",
	}

	if err := SendJSON(&buf, msg); err != nil {
		t.Fatalf("failed to send JSON: %v", err)
	}

	readBack, err := ReadMessage(&buf)
	if err != nil {
		t.Fatalf("failed to read back JSON message: %v", err)
	}

	expected := `{"event":"handshake","data":"ok"}`
	if string(readBack) != expected {
		t.Fatalf("expected JSON %q, got %q", expected, string(readBack))
	}
}

func TestGetDefaultSocketPath(t *testing.T) {
	// Test environment variable override
	customPath := "/tmp/custom_nagare.sock"
	t.Setenv("NAGARE_PLUGIN_SOCKET_PATH", customPath)

	path := GetDefaultSocketPath()
	if path != customPath {
		t.Fatalf("expected env overridden socket path %q, got %q", customPath, path)
	}

	// Test default path when env var is unset
	_ = os.Unsetenv("NAGARE_PLUGIN_SOCKET_PATH")
	defaultPath := GetDefaultSocketPath()
	if defaultPath == "" {
		t.Fatalf("expected non-empty default socket path")
	}
}
