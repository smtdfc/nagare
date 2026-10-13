package provider

import (
	"context"
	"testing"

	"github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type mockAdapter struct{}

func (m *mockAdapter) GetModels(context.Context) ([]string, error) {
	return []string{"mock-model-1"}, nil
}

func (m *mockAdapter) Send(context.Context, string, messages.ListMessage, tool.ListTool) (message.ReadOnlyChannel, error) {
	ch := make(chan messages.Message)
	close(ch)
	return ch, nil
}

func (m *mockAdapter) Embed(_ context.Context, _ string, inputs []string) ([][]float32, error) {
	result := make([][]float32, len(inputs))
	for i := range inputs {
		result[i] = []float32{float32(i)}
	}
	return result, nil
}

func TestModelProviderAdapter_Interface(t *testing.T) {
	// Verify mock adapter satisfies ModelProviderAdapter interface
	var adapter ModelProviderAdapter = &mockAdapter{}

	models, err := adapter.GetModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 1 || models[0] != "mock-model-1" {
		t.Errorf("expected ['mock-model-1'], got %v", models)
	}

	ch, err := adapter.Send(context.Background(), "mock-model-1", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error on Send: %v", err)
	}
	if ch == nil {
		t.Errorf("expected non-nil channel")
	}
}
