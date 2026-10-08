package client

import (
	"testing"
)

func TestConnectConfigStructure(t *testing.T) {
	cfg := ConnectConfig{
		SocketPath:  "/tmp/nagare.sock",
		ConnectCode: "secret-code-123",
	}

	if cfg.SocketPath != "/tmp/nagare.sock" || cfg.ConnectCode != "secret-code-123" {
		t.Fatalf("unexpected values in ConnectConfig: %+v", cfg)
	}
}
