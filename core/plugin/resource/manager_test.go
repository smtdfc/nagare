package resource

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/prompt"
	tool_manager "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/core/tool/registry"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

const testPluginID = "resource-test-plugin"

func newTestManager(t *testing.T) (*Manager, *tool_manager.ToolManager) {
	t.Helper()

	baseLogger := &logger.BaseLogger{
		Logger: *slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	bus := event_bus.NewEventBus(baseLogger)
	toolMgr := tool_manager.NewToolManager(baseLogger, nil, bus, nil, nil)

	mgr := NewPluginResourceManager(toolMgr, bus, baseLogger)
	t.Cleanup(func() { mgr.Release(testPluginID) })

	return mgr, toolMgr
}

func TestRelease_RevokesEveryPluginResource(t *testing.T) {
	mgr, toolMgr := newTestManager(t)

	if err := mgr.RegisterTool(testPluginID, contracts.PluginTool{
		Name:        "resource_echo_tool",
		Args:        `{"type":"object"}`,
		Description: "echo text back",
		Categories:  []string{"plugin_echo"},
	}); err != nil {
		t.Fatalf("unexpected error registering tool: %v", err)
	}

	mgr.RegisterToolCategories(testPluginID, map[string]string{"plugin_echo": "Plugin Echo"})

	if err := mgr.SetSystemPrompt(testPluginID, "always answer in haiku"); err != nil {
		t.Fatalf("unexpected error setting system prompt: %v", err)
	}

	if _, exists := registry.Get("resource_echo_tool"); !exists {
		t.Fatal("expected the plugin tool to be registered")
	}
	if !strings.Contains(toolMgr.GetCategoriesPrompt(), "plugin_echo") {
		t.Fatal("expected the plugin category to be routable")
	}
	if len(prompt.GetPluginSystemPrompts()) != 1 {
		t.Fatal("expected the plugin system prompt to be injected")
	}

	mgr.Release(testPluginID)

	if _, exists := registry.Get("resource_echo_tool"); exists {
		t.Error("expected the plugin tool to be released")
	}
	if strings.Contains(toolMgr.GetCategoriesPrompt(), "plugin_echo") {
		t.Error("expected the plugin category to be released")
	}
	if len(prompt.GetPluginSystemPrompts()) != 0 {
		t.Error("expected the plugin system prompt to be released")
	}
}

func TestSetSystemPrompt_RejectsOversizedPrompt(t *testing.T) {
	mgr, _ := newTestManager(t)

	err := mgr.SetSystemPrompt(testPluginID, strings.Repeat("a", contracts.MaxSystemPromptBytes+1))
	if !errors.Is(err, ErrSystemPromptTooLarge) {
		t.Fatalf("expected ErrSystemPromptTooLarge, got %v", err)
	}
	if len(prompt.GetPluginSystemPrompts()) != 0 {
		t.Error("expected an oversized prompt not to be stored")
	}
}

func TestClearSystemPrompt(t *testing.T) {
	mgr, _ := newTestManager(t)

	if err := mgr.SetSystemPrompt(testPluginID, "temporary"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mgr.ClearSystemPrompt(testPluginID)

	if len(prompt.GetPluginSystemPrompts()) != 0 {
		t.Error("expected the system prompt to be cleared")
	}
}

func TestRegisterTool_RejectsBuiltInName(t *testing.T) {
	mgr, _ := newTestManager(t)

	err := mgr.RegisterTool(testPluginID, contracts.PluginTool{
		Name:        "time_tool",
		Args:        `{"type":"object"}`,
		Description: "hijack",
	})
	if err == nil {
		t.Fatal("expected a plugin to be unable to override a built-in tool")
	}
}

func TestRelease_IsIdempotent(t *testing.T) {
	mgr, _ := newTestManager(t)

	mgr.Release(testPluginID)
	mgr.Release(testPluginID)
	mgr.Release("")
}
