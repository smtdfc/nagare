package host

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/smtdfc/nagare/core/prompt"
	"github.com/smtdfc/nagare/pkgs/ipc"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func readResponse(t *testing.T, conn net.Conn) contracts.Payload[map[string]any] {
	t.Helper()

	raw, err := ipc.ReadMessage(conn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	var payload contracts.Payload[map[string]any]
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	return payload
}

func TestHandleInjectSystemPrompt_RequiresScope(t *testing.T) {
	host := newTestHost()
	server, client := net.Pipe()
	t.Cleanup(func() { _ = client.Close() })
	t.Cleanup(func() { prompt.ClearPluginSystemPrompt("plugin-1") })

	conn := &PluginConnection{
		conn:     server,
		pluginID: "plugin-1",
		scopes:   []string{"chat"},
	}
	payload := &contracts.Payload[any]{
		RequestID: "request-1",
		Data: map[string]any{
			"content": "injected context",
		},
	}

	go host.handleInjectSystemPrompt(conn, payload)

	resp := readResponse(t, client)
	if resp.Event != contracts.InjectSystemPromptFailedEvent {
		t.Fatalf("expected failed event, got %s", resp.Event)
	}
	if got := prompt.GetPluginSystemPrompts(); len(got) != 0 {
		t.Fatalf("prompt must not be stored without the system_prompt scope, got %+v", got)
	}
}

func TestHandleInjectAndClearSystemPrompt(t *testing.T) {
	host := newTestHost()
	t.Cleanup(func() { prompt.ClearPluginSystemPrompt("plugin-1") })

	server, client := net.Pipe()
	t.Cleanup(func() { _ = client.Close() })

	conn := &PluginConnection{
		conn:     server,
		pluginID: "plugin-1",
		scopes:   []string{"system_prompt"},
	}

	go host.handleInjectSystemPrompt(conn, &contracts.Payload[any]{
		RequestID: "request-1",
		Data: map[string]any{
			"content": "injected context",
		},
	})

	if resp := readResponse(t, client); resp.Event != contracts.InjectSystemPromptSuccessEvent {
		t.Fatalf("expected inject success event, got %s", resp.Event)
	}
	prompts := prompt.GetPluginSystemPrompts()
	if len(prompts) != 1 || prompts[0].PluginID != "plugin-1" || prompts[0].Content != "injected context" {
		t.Fatalf("expected injected prompt, got %+v", prompts)
	}

	go host.handleClearSystemPrompt(conn, &contracts.Payload[any]{
		RequestID: "request-2",
	})

	if resp := readResponse(t, client); resp.Event != contracts.ClearSystemPromptSuccessEvent {
		t.Fatalf("expected clear success event, got %s", resp.Event)
	}
	if prompts := prompt.GetPluginSystemPrompts(); len(prompts) != 0 {
		t.Fatalf("expected cleared prompt, got %+v", prompts)
	}
}
