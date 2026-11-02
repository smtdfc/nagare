package host

import (
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
)

func newTestLogger() *logger.BaseLogger {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	return &logger.BaseLogger{Logger: *slog.New(handler)}
}

func newTestHost() *PluginHost {
	return &PluginHost{
		logger:      newTestLogger(),
		eventBus:    event_bus.NewEventBus(newTestLogger()),
		connections: make(map[string]*PluginConnection),
		stopChan:    make(chan struct{}),
	}
}

// A stale connection (rejected during handshake, or from before a restart) must
// never unregister the plugin that currently owns the id.
func TestUnregister_KeepsNewOwner(t *testing.T) {
	host := newTestHost()

	current := &PluginConnection{pluginID: "plugin-1"}
	stale := &PluginConnection{pluginID: "plugin-1"}

	host.connections[current.pluginID] = current

	host.unregister(stale)

	if host.connections[current.pluginID] != current {
		t.Errorf("stale unregister removed the active connection")
	}
}

func TestUnregister_ReleasesOwner(t *testing.T) {
	host := newTestHost()

	conn := &PluginConnection{pluginID: "plugin-1"}
	host.connections[conn.pluginID] = conn

	host.unregister(conn)

	if len(host.connections) != 0 {
		t.Errorf("expected connections to be empty, got %d", len(host.connections))
	}
}

// Stop used to take h.mu and then call unregister (which takes h.mu again),
// deadlocking as soon as one plugin was connected.
func TestStop_NoDeadlockAndClosesConnections(t *testing.T) {
	host := newTestHost()

	server, client := net.Pipe()
	t.Cleanup(func() { _ = client.Close() })

	conn := &PluginConnection{conn: server, pluginID: "plugin-1"}
	host.connections[conn.pluginID] = conn

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = host.Stop()
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop deadlocked with an active connection")
	}

	if len(host.connections) != 0 {
		t.Errorf("expected connections to be cleared, got %d", len(host.connections))
	}

	_ = client.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := client.Read(make([]byte, 1)); err == nil {
		t.Errorf("expected client side of the pipe to be closed")
	}
}
