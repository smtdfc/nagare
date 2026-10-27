package host

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/internal"
	"github.com/smtdfc/nagare/pkgs/ipc"
)

func newTestLogger() *logger.BaseLogger {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	return &logger.BaseLogger{Logger: *slog.New(handler)}
}

func newTestHost() *PluginHost {
	return &PluginHost{
		logger:        newTestLogger(),
		eventBus:      event_bus.NewEventBus(newTestLogger()),
		connections:   make(map[string]*PluginConnection),
		specialGroups: make(map[PluginGroup]*PluginConnection),
		stopChan:      make(chan struct{}),
	}
}

func TestGetGroupName(t *testing.T) {
	tests := []struct {
		name     string
		features []plugin.Feature
		want     PluginGroup
	}{
		{name: "no features", features: nil, want: NormalPluginGroup},
		{name: "chat", features: []plugin.Feature{plugin.ChatFeature}, want: NormalPluginGroup},
		{
			name:     "vector index",
			features: []plugin.Feature{plugin.VectorIndexFeature},
			want:     VectorIndexPluginGroup,
		},
		{
			name:     "vector index with other features",
			features: []plugin.Feature{plugin.ChatFeature, plugin.VectorIndexFeature},
			want:     VectorIndexPluginGroup,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := getGroupName(tc.features); got != tc.want {
				t.Errorf("getGroupName(%v) = %q, want %q", tc.features, got, tc.want)
			}
		})
	}
}

// A stale connection (rejected during handshake, or from before a restart) must
// never unregister the plugin that currently owns the id/group.
func TestUnregister_KeepsNewOwner(t *testing.T) {
	host := newTestHost()

	current := &PluginConnection{pluginID: "plugin-1"}
	stale := &PluginConnection{pluginID: "plugin-1"}

	host.connections[current.pluginID] = current
	host.specialGroups[VectorIndexPluginGroup] = current

	host.unregister(stale)

	if host.connections[current.pluginID] != current {
		t.Errorf("stale unregister removed the active connection")
	}
	if host.specialGroups[VectorIndexPluginGroup] != current {
		t.Errorf("stale unregister removed the active special group entry")
	}
}

func TestUnregister_ReleasesOwner(t *testing.T) {
	host := newTestHost()

	conn := &PluginConnection{pluginID: "plugin-1"}
	host.connections[conn.pluginID] = conn
	host.specialGroups[VectorIndexPluginGroup] = conn

	host.unregister(conn)

	if len(host.connections) != 0 {
		t.Errorf("expected connections to be empty, got %d", len(host.connections))
	}
	if len(host.specialGroups) != 0 {
		t.Errorf("expected specialGroups to be empty, got %d", len(host.specialGroups))
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
	host.specialGroups[VectorIndexPluginGroup] = conn

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
	if len(host.specialGroups) != 0 {
		t.Errorf("expected specialGroups to be cleared, got %d", len(host.specialGroups))
	}

	_ = client.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := client.Read(make([]byte, 1)); err == nil {
		t.Errorf("expected client side of the pipe to be closed")
	}
}

// The forwarder must survive an add event that arrives before the vector plugin
// is installed, and must forward the next event once the plugin connects.
func TestForwardVectorIndexActions_SurvivesMissingPlugin(t *testing.T) {
	host := newTestHost()
	go host.forwardVectorIndexActions()
	t.Cleanup(func() { close(host.stopChan) })

	errCh, unsubscribeErr := host.eventBus.Subscribe(event_bus.AddVectorErrorEvent)
	t.Cleanup(unsubscribeErr)

	// The forwarder subscribes on its own goroutine, so publish until it is
	// listening instead of racing the goroutine startup.
	deadline := time.Now().Add(2 * time.Second)
	receivedError := false
	for !receivedError && time.Now().Before(deadline) {
		host.eventBus.Publish(context.Background(), event_bus.AddVectorEvent, &event_bus.AddVectorEventPayload{
			RequestID: "req-missing",
			Vector:    []float32{1, 2, 3, 4},
			IDs:       []uint64{1},
		})

		select {
		case evt := <-errCh:
			payload, ok := evt.(*event_bus.AddVectorErrorEventPayload)
			if !ok {
				t.Fatalf("unexpected error payload type %T", evt)
			}
			if payload.RequestID != "req-missing" {
				t.Errorf("expected requestID %q, got %q", "req-missing", payload.RequestID)
			}
			if payload.Error == nil {
				t.Errorf("expected a non-nil error when the vector plugin is missing")
			}
			receivedError = true
		case <-time.After(20 * time.Millisecond):
		}
	}
	if !receivedError {
		t.Fatal("expected an error event when the vector plugin is missing")
	}

	server, client := net.Pipe()
	t.Cleanup(func() {
		_ = server.Close()
		_ = client.Close()
	})

	type envelope struct {
		Event     string `json:"event"`
		RequestID string `json:"requestID"`
	}

	received := make(chan envelope, 64)
	go func() {
		defer close(received)
		for {
			msg, err := ipc.ReadMessage(client)
			if err != nil {
				return
			}
			var env envelope
			if err := json.Unmarshal(msg, &env); err == nil {
				received <- env
			}
		}
	}()

	conn := &PluginConnection{conn: server, pluginID: "plugin-1"}
	host.mu.Lock()
	host.specialGroups[VectorIndexPluginGroup] = conn
	host.mu.Unlock()

	host.eventBus.Publish(context.Background(), event_bus.AddVectorEvent, &event_bus.AddVectorEventPayload{
		RequestID: "req-forwarded",
		Vector:    []float32{1, 2, 3, 4},
		IDs:       []uint64{1},
	})

	timeout := time.After(2 * time.Second)
	for {
		select {
		case env, ok := <-received:
			if !ok {
				t.Fatal("plugin connection closed before the event was forwarded")
			}
			if env.RequestID != "req-forwarded" {
				continue // an event published before the plugin connected
			}
			if env.Event != string(internal.AddVectorEvent) {
				t.Errorf("expected event %q, got %q", internal.AddVectorEvent, env.Event)
			}
			return
		case <-timeout:
			t.Fatal("vector add event was not forwarded after the plugin connected")
		}
	}
}
