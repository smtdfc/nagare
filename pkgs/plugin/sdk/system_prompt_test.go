package sdk

import (
	"context"
	"errors"
	"strings"
	"testing"

	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func TestSystemPromptClientOperationsWithoutConnection(t *testing.T) {
	p := NewPlugin()
	ctx := context.Background()

	if err := p.InjectSystemPrompt(ctx, "extra context"); err == nil {
		t.Fatalf("expected error when injecting a system prompt without a connection")
	}

	if err := p.ClearSystemPrompt(ctx); err == nil {
		t.Fatalf("expected error when clearing a system prompt without a connection")
	}
}

func TestInjectSystemPrompt_RejectsOversizedContent(t *testing.T) {
	p := NewPlugin()
	err := p.InjectSystemPrompt(
		context.Background(),
		strings.Repeat("a", contracts.MaxSystemPromptBytes+1),
	)

	if !errors.Is(err, ErrSystemPromptTooLarge) {
		t.Fatalf("expected ErrSystemPromptTooLarge, got %v", err)
	}
}
