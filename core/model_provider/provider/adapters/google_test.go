package adapters

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/messages"
	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/apierrors"
	"google.golang.org/genai/interactions/models/components"
	"google.golang.org/genai/interactions/models/interactions"
)

func newTestGoogleAdapter() *GoogleAdapter {
	return &GoogleAdapter{
		Models: []string{"gemini-flash-latest"},
		logger: testLogger(),
	}
}

// TestTransformToProviderInput_MapsRoles verifies system messages are dropped
// from the step list (they belong in SystemInstruction) and that user/agent
// text map to the right step type.
func TestTransformToProviderInput_MapsRoles(t *testing.T) {
	adapter := newTestGoogleAdapter()

	steps, err := adapter.TransformToProviderInput(messages.ListMessage{
		messages.NewTextMessage(messages.SYSTEM, "be brief"),
		messages.NewTextMessage(messages.USER, "hello"),
		messages.NewTextMessage(messages.AGENT, "hi"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(steps))
	}
	if steps[0].UserInputStep == nil {
		t.Fatal("expected first step to be a user input step")
	}
	if steps[1].ModelOutputStep == nil {
		t.Fatal("expected second step to be a model output step")
	}
}

// TestTransformToSystemInstruction verifies system and developer messages are
// joined into a single instruction string.
func TestTransformToSystemInstruction(t *testing.T) {
	adapter := newTestGoogleAdapter()

	instruction := adapter.TransformToSystemInstruction(messages.ListMessage{
		messages.NewTextMessage(messages.SYSTEM, "rule one"),
		messages.NewTextMessage(messages.DEVELOPER, "rule two"),
		messages.NewTextMessage(messages.USER, "ignore me"),
	})

	if instruction == nil {
		t.Fatal("expected a system instruction")
	}
	if *instruction != "rule one\n\nrule two" {
		t.Errorf("unexpected instruction: %q", *instruction)
	}
}

// TestTransformToProviderInput_ToolCall verifies tool call and result messages
// are mapped with matching IDs and parsed arguments.
func TestTransformToProviderInput_ToolCall(t *testing.T) {
	adapter := newTestGoogleAdapter()

	steps, err := adapter.TransformToProviderInput(messages.ListMessage{
		messages.NewToolCallMessage("call-1", "read_file", `{"path":"a.txt"}`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(steps) != 1 || steps[0].FunctionCallStep == nil {
		t.Fatalf("expected a function call step, got %#v", steps)
	}

	call := steps[0].FunctionCallStep
	if call.ID != "call-1" || call.Name != "read_file" {
		t.Errorf("unexpected call: %#v", call)
	}
	if path, ok := call.Arguments["path"]; !ok || path != "a.txt" {
		t.Errorf("unexpected arguments: %#v", call.Arguments)
	}

	resultStep, err := adapter.TransformToProviderMessage(&messages.ToolResultMessage{
		CallID: "call-1",
		Name:   "read_file",
		Result: "content",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resultStep.FunctionResultStep == nil {
		t.Fatal("expected a function result step")
	}
	if resultStep.FunctionResultStep.CallID != "call-1" {
		t.Errorf("unexpected call id: %q", resultStep.FunctionResultStep.CallID)
	}
}

// TestBuildCreateRequest_IsStateless is the guard for the whole design: the
// request must opt out of server-side storage and must never reference a
// previous interaction.
func TestBuildCreateRequest_IsStateless(t *testing.T) {
	adapter := newTestGoogleAdapter()

	body, err := adapter.buildCreateRequest(
		"gemini-flash-latest",
		messages.ListMessage{
			messages.NewTextMessage(messages.SYSTEM, "be brief"),
			messages.NewTextMessage(messages.USER, "hello"),
		},
		tool.ListTool{},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	payload := string(encoded)
	if !strings.Contains(payload, `"store":false`) {
		t.Errorf("expected store to be false, got %s", payload)
	}
	if strings.Contains(payload, "previous_interaction_id") {
		t.Errorf("expected no previous interaction reference, got %s", payload)
	}
	if strings.Contains(payload, `"stream":false`) {
		t.Errorf("expected streaming to be enabled, got %s", payload)
	}
}

// TestMapCreateError_ServiceUnavailable covers the real 503 body shape: the
// SDK returns CreateInteractionServerError and the cause must be the human
// readable message, not the raw JSON.
func TestMapCreateError_ServiceUnavailable(t *testing.T) {
	adapter := newTestGoogleAdapter()

	err := &apierrors.CreateInteractionServerError{
		Error_: interactions.Error{
			Code:    genai.Ptr("service_unavailable"),
			Message: genai.Ptr("gemini-3.8-flash is currently experiencing high demand"),
		},
		HTTPMeta: components.HTTPMetadata{Response: &http.Response{StatusCode: 503}},
	}

	code, cause := adapter.mapCreateError(err, "gemini-3.8-flash")
	if code != "500" {
		t.Errorf("expected code 500, got %q", code)
	}
	if !strings.Contains(cause, "high demand") {
		t.Errorf("expected the provider message, got %q", cause)
	}
	if strings.Contains(cause, "{") {
		t.Errorf("expected no raw JSON, got %q", cause)
	}
}

// TestTransformToProviderMessage_Unsupported verifies unsupported message
// types fail loudly instead of producing an empty step, which the API cannot
// serialize.
func TestTransformToProviderMessage_Unsupported(t *testing.T) {
	adapter := newTestGoogleAdapter()

	if _, err := adapter.TransformToProviderMessage(messages.NewReasoningMessage("thinking")); err == nil {
		t.Fatal("expected an error for an unsupported message type")
	}
}
