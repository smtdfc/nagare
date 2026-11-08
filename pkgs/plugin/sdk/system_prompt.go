package sdk

import (
	"context"
	"errors"
	"fmt"

	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

// InjectSystemPrompt appends plugin-owned context to Nagare's base system
// prompt for every chat invoke, regardless of session or channel. Re-injecting
// replaces the plugin's previous fragment. Passing empty content clears the
// fragment instead of adding one.
//
// Injected prompts live in the gateway process. A plugin must inject them
// again after reconnecting.
func (p *PluginClient) InjectSystemPrompt(ctx context.Context, content string) error {
	if len(content) > contracts.MaxSystemPromptBytes {
		return fmt.Errorf("%w: maximum is %d bytes", ErrSystemPromptTooLarge, contracts.MaxSystemPromptBytes)
	}

	resp, err := sendAndWait[contracts.InjectSystemPromptSuccessEventPayload, contracts.InjectSystemPromptFailedEventPayload](
		p,
		ctx,
		contracts.InjectSystemPromptEvent,
		contracts.InjectSystemPromptEventPayload{
			Content: content,
		},
		contracts.InjectSystemPromptFailedEvent,
		contracts.InjectSystemPromptSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to inject system prompt", "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to inject system prompt", "error", resp.Error)
		return errors.New(resp.Error.Cause)
	}

	return nil
}

// ClearSystemPrompt removes the system prompt fragment injected with
// InjectSystemPrompt.
func (p *PluginClient) ClearSystemPrompt(ctx context.Context) error {
	resp, err := sendAndWait[contracts.ClearSystemPromptSuccessEventPayload, contracts.ClearSystemPromptFailedEventPayload](
		p,
		ctx,
		contracts.ClearSystemPromptEvent,
		contracts.ClearSystemPromptEventPayload{},
		contracts.ClearSystemPromptFailedEvent,
		contracts.ClearSystemPromptSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to clear system prompt", "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to clear system prompt", "error", resp.Error)
		return errors.New(resp.Error.Cause)
	}

	return nil
}
