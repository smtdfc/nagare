package session

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/pkgs/messages"
)

func TestSessionState_Construction(t *testing.T) {
	// Verify SessionState initialization and field retention
	sessID := uuid.New()
	providerID := uuid.New()

	state := SessionState{
		SessionID:          sessID,
		CurrentModel:       "gpt-4o-mini",
		CurrentLLMProvider: providerID,
		ChannelID:          "ch-abc",
		Messages:           messages.ListMessage{},
		NextCursor:         "cursor-1",
		OwnerType:          USER,
		OwnerID:            "u-1",
	}

	if state.SessionID != sessID {
		t.Errorf("expected session ID %s, got %s", sessID, state.SessionID)
	}
	if state.CurrentModel != "gpt-4o-mini" {
		t.Errorf("expected model gpt-4o-mini, got %s", state.CurrentModel)
	}
	if state.ChannelID != "ch-abc" {
		t.Errorf("expected channel ch-abc, got %s", state.ChannelID)
	}
}
