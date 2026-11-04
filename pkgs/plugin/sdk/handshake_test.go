package sdk

import (
	"context"
	"testing"

	"github.com/smtdfc/nagare/pkgs/plugin/metadata"
)

func TestHandshakeWithoutConnection(t *testing.T) {
	p := NewPlugin()
	p.Metadata = &metadata.PluginMetadata{
		PackageName: "nagare.plugin.test",
	}

	// Should return error because connector has no active connection
	err := p.Handshake(context.Background())
	if err == nil {
		t.Fatalf("expected error when handshaking without connection, got nil")
	}
}
