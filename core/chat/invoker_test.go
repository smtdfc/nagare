package chat

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/config"
)

func TestAgentInvoker_NilParams(t *testing.T) {
	// Verify error when params is nil
	invoker := &AgentInvoker{}
	_, err := invoker.Invoke(context.Background(), nil)
	if err == nil {
		t.Errorf("expected error when invoke params is nil")
	}
}

func TestSelectCompactModel(t *testing.T) {
	sessionProvider := uuid.New()
	summaryProvider := uuid.New()

	tests := []struct {
		name          string
		generalConfig *config.GeneralConfig
		wantProvider  uuid.UUID
		wantModel     string
		wantErr       bool
	}{
		{
			name:          "nil config falls back to session model",
			generalConfig: nil,
			wantProvider:  sessionProvider,
			wantModel:     "session-model",
		},
		{
			name:          "empty summary config falls back to session model",
			generalConfig: &config.GeneralConfig{},
			wantProvider:  sessionProvider,
			wantModel:     "session-model",
		},
		{
			name: "summary model takes precedence",
			generalConfig: &config.GeneralConfig{
				SummaryModelProvider: summaryProvider.String(),
				SummaryModel:         "summary-model",
			},
			wantProvider: summaryProvider,
			wantModel:    "summary-model",
		},
		{
			name: "invalid summary provider falls back with error",
			generalConfig: &config.GeneralConfig{
				SummaryModelProvider: "not-a-uuid",
				SummaryModel:         "summary-model",
			},
			wantProvider: uuid.Nil,
			wantModel:    "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, model, err := selectCompactModel(sessionProvider, "session-model", tt.generalConfig)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error state: %v", err)
			}
			if provider != tt.wantProvider {
				t.Errorf("expected provider %s, got %s", tt.wantProvider, provider)
			}
			if model != tt.wantModel {
				t.Errorf("expected model %q, got %q", tt.wantModel, model)
			}
		})
	}
}

func TestAgentInvokeParams_Construction(t *testing.T) {
	// Verify AgentInvokeParams fields
	params := &AgentInvokeParams{
		SessionID:        "sess-123",
		SenderType:       "user",
		SenderID:         "u-1",
		SendIntoEventBus: true,
		InvokeID:         "inv-1",
	}

	if params.SessionID != "sess-123" {
		t.Errorf("expected session ID 'sess-123', got '%s'", params.SessionID)
	}
	if !params.SendIntoEventBus {
		t.Errorf("expected SendIntoEventBus true")
	}
}
