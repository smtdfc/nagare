package adapters

import (
	"io"
	"log/slog"
	"testing"

	"github.com/openai/openai-go/v3/responses"
	nagare_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/messages"
)

func testLogger() *logger.BaseLogger {
	return &logger.BaseLogger{Logger: *slog.New(slog.NewTextHandler(io.Discard, nil))}
}

func newTestAdapter() *OpenAICompatibleAdapter {
	return &OpenAICompatibleAdapter{
		BaseURL: "http://localhost",
		Models:  []string{"gpt-4"},
		logger:  testLogger(),
	}
}

// TestTransformToProviderMessage_TextMessage_User verifies USER role mapping.
func TestTransformToProviderMessage_TextMessage_User(t *testing.T) {
	adapter := newTestAdapter()
	msg := messages.NewTextMessage(messages.USER, "hello")

	result, err := adapter.TransformToProviderMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OfInputMessage == nil {
		t.Fatal("expected OfInputMessage to be set")
	}
	if result.OfInputMessage.Role != "user" {
		t.Errorf("expected role 'user', got %q", result.OfInputMessage.Role)
	}
}

// TestTransformToProviderMessage_TextMessage_System verifies SYSTEM role mapping.
func TestTransformToProviderMessage_TextMessage_System(t *testing.T) {
	adapter := newTestAdapter()
	msg := messages.NewTextMessage(messages.SYSTEM, "system prompt")

	result, err := adapter.TransformToProviderMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OfInputMessage == nil {
		t.Fatal("expected OfInputMessage to be set")
	}
	if result.OfInputMessage.Role != "system" {
		t.Errorf("expected role 'system', got %q", result.OfInputMessage.Role)
	}
}

// TestTransformToProviderMessage_TextMessage_Developer verifies DEVELOPER role mapping.
func TestTransformToProviderMessage_TextMessage_Developer(t *testing.T) {
	adapter := newTestAdapter()
	msg := messages.NewTextMessage(messages.DEVELOPER, "dev message")

	result, err := adapter.TransformToProviderMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OfInputMessage == nil {
		t.Fatal("expected OfInputMessage to be set")
	}
	if result.OfInputMessage.Role != "developer" {
		t.Errorf("expected role 'developer', got %q", result.OfInputMessage.Role)
	}
}

// TestTransformToProviderMessage_TextMessage_Agent verifies AGENT role uses OfOutputMessage.
func TestTransformToProviderMessage_TextMessage_Agent(t *testing.T) {
	adapter := newTestAdapter()
	msg := messages.NewTextMessage(messages.AGENT, "agent response")

	result, err := adapter.TransformToProviderMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OfOutputMessage == nil {
		t.Fatal("expected OfOutputMessage to be set for AGENT role")
	}
}

// TestTransformToProviderMessage_ToolCallMessage verifies tool call mapping.
func TestTransformToProviderMessage_ToolCallMessage(t *testing.T) {
	adapter := newTestAdapter()
	msg := messages.NewToolCallMessage("call-1", "test_tool", `{"arg":"value"}`)

	result, err := adapter.TransformToProviderMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OfFunctionCall == nil {
		t.Fatal("expected OfFunctionCall to be set")
	}
	if result.OfFunctionCall.CallID != "call-1" {
		t.Errorf("expected CallID 'call-1', got %q", result.OfFunctionCall.CallID)
	}
	if result.OfFunctionCall.Name != "test_tool" {
		t.Errorf("expected Name 'test_tool', got %q", result.OfFunctionCall.Name)
	}
}

// TestTransformToProviderMessage_ToolResultMessage verifies tool result mapping.
func TestTransformToProviderMessage_ToolResultMessage(t *testing.T) {
	adapter := newTestAdapter()
	msg := messages.NewToolResultMessage("call-1", "test_tool", `{"result":"ok"}`)

	result, err := adapter.TransformToProviderMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OfFunctionCallOutput == nil {
		t.Fatal("expected OfFunctionCallOutput to be set")
	}
	if result.OfFunctionCallOutput.CallID != "call-1" {
		t.Errorf("expected CallID 'call-1', got %q", result.OfFunctionCallOutput.CallID)
	}
}

// TestTransformToolDeclarations_Empty verifies empty tool list returns empty result.
func TestTransformToolDeclarations_Empty(t *testing.T) {
	adapter := newTestAdapter()
	result, err := adapter.TransformToolDeclarations(tool.ListTool{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 tools, got %d", len(result))
	}
}

// TestTransformToolDeclarations_ValidTool verifies a tool with valid JSON schema is transformed.
func TestTransformToolDeclarations_ValidTool(t *testing.T) {
	adapter := newTestAdapter()

	mockTool := &mockToolForAdapter{
		name:        "test_tool",
		description: "A test tool",
		argsSchema:  `{"type":"object","properties":{"name":{"type":"string"}}}`,
	}

	result, err := adapter.TransformToolDeclarations(tool.ListTool{mockTool})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(result))
	}
	if result[0].OfFunction == nil {
		t.Fatal("expected OfFunction to be set")
	}
	if result[0].OfFunction.Name != "test_tool" {
		t.Errorf("expected name 'test_tool', got %q", result[0].OfFunction.Name)
	}
}

// TestTransformToolDeclarations_InvalidSchema verifies invalid JSON schema returns an error.
func TestTransformToolDeclarations_InvalidSchema(t *testing.T) {
	adapter := newTestAdapter()

	mockTool := &mockToolForAdapter{
		name:        "bad_tool",
		description: "Bad schema",
		argsSchema:  `{invalid json`,
	}

	_, err := adapter.TransformToolDeclarations(tool.ListTool{mockTool})
	if err == nil {
		t.Fatal("expected error for invalid schema, got nil")
	}
}

// TestTransformToProviderMessage_UnknownType verifies unknown message types return zero value.
func TestTransformToProviderMessage_UnknownType(t *testing.T) {
	adapter := newTestAdapter()

	// Use a custom message type that doesn't match any case
	msg := &unknownMessage{}

	result, err := adapter.TransformToProviderMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should return zero value (all nil)
	empty := responses.ResponseInputItemUnionParam{}
	if result != empty {
		t.Error("expected zero-value ResponseInputItemUnionParam for unknown message type")
	}
}

// mockToolForAdapter implements tool.Tool for testing TransformToolDeclarations.
type mockToolForAdapter struct {
	name        string
	description string
	argsSchema  string
}

func (m *mockToolForAdapter) GetName() string               { return m.name }
func (m *mockToolForAdapter) GetDescription() string         { return m.description }
func (m *mockToolForAdapter) GetArgsSchema() string          { return m.argsSchema }
func (m *mockToolForAdapter) GetBindings() tool.Bindings     { return nil }
func (m *mockToolForAdapter) GetCategories() []string        { return nil }
func (m *mockToolForAdapter) WithBindings(_ tool.Bindings) tool.Tool { return m }
func (m *mockToolForAdapter) Execute(_ *nagare_context.ExecuteContext, _ string) (string, error) {
	return "", nil
}

// unknownMessage is a message type not handled by TransformToProviderMessage.
type unknownMessage struct{}

func (u *unknownMessage) GetMessageID() string           { return "unknown" }
func (u *unknownMessage) GetMessageType() messages.MessageType { return "unknown" }
func (u *unknownMessage) GetInvokeID() string            { return "" }
func (u *unknownMessage) SetInvokeID(_ string)           {}
