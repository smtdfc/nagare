package host

import (
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/plugin/resource"
	"github.com/smtdfc/nagare/core/prompt"
	tool_manager "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/core/tool/registry"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func newTestLogger() *logger.BaseLogger {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	return &logger.BaseLogger{Logger: *slog.New(handler)}
}

func newTestHost() *PluginHost {
	baseLogger := newTestLogger()
	eventBus := event_bus.NewEventBus(baseLogger)
	toolMgr := tool_manager.NewToolManager(baseLogger, nil, eventBus, nil, nil)

	return &PluginHost{
		logger:        baseLogger,
		eventBus:      eventBus,
		resourceMgr:   resource.NewPluginResourceManager(toolMgr, eventBus, baseLogger),
		connections:   make(map[string]*PluginConnection),
		chatListeners: make(map[string]*PluginConnection),
		stopChan:      make(chan struct{}),
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

func registerTestPluginResources(t *testing.T, host *PluginHost, pluginID string) {
	t.Helper()

	if err := host.resourceMgr.RegisterTool(pluginID, contracts.PluginTool{
		Name:        "host_test_tool",
		Args:        `{"type":"object"}`,
		Description: "test tool",
		Categories:  []string{"host_test"},
	}); err != nil {
		t.Fatalf("unexpected error registering plugin tool: %v", err)
	}

	host.resourceMgr.RegisterToolCategories(pluginID, map[string]string{"host_test": "Host Test"})

	if err := host.resourceMgr.SetSystemPrompt(pluginID, "host test prompt"); err != nil {
		t.Fatalf("unexpected error setting plugin system prompt: %v", err)
	}
}

// A plugin process that disconnects must not leave its injected tools, routing
// categories or system prompt behind.
func TestUnregister_ReleasesPluginResources(t *testing.T) {
	host := newTestHost()
	t.Cleanup(func() { registry.UnregisterPluginTools("plugin-release") })

	conn := &PluginConnection{pluginID: "plugin-release"}
	host.connections[conn.pluginID] = conn
	registerTestPluginResources(t, host, conn.pluginID)

	host.unregister(conn)

	if _, exists := registry.Get("host_test_tool"); exists {
		t.Error("expected the plugin tool to be released on disconnect")
	}
	if len(prompt.GetPluginSystemPrompts()) != 0 {
		t.Error("expected the plugin system prompt to be released on disconnect")
	}
}

// A stale connection (for example one rejected during handshake) must never
// release resources that belong to the current owner of the plugin id.
func TestUnregister_StaleConnectionKeepsResources(t *testing.T) {
	host := newTestHost()
	t.Cleanup(func() { registry.UnregisterPluginTools("plugin-stale") })

	current := &PluginConnection{pluginID: "plugin-stale"}
	stale := &PluginConnection{pluginID: "plugin-stale"}
	host.connections[current.pluginID] = current
	registerTestPluginResources(t, host, current.pluginID)

	host.unregister(stale)

	if _, exists := registry.Get("host_test_tool"); !exists {
		t.Error("expected a stale connection to leave plugin resources untouched")
	}

	host.unregister(current)
}
