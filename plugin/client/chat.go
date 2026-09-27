package client

import (
	"context"
	"errors"

	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
)

func (p *PluginClient) PrepareChatSession(ctx context.Context, channelID string) (string, error) {
	payload := plugin_dtos.PrepareChatSessionEventPayload{
		ChannelID: channelID,
	}

	resp, err := sendAndWait[plugin_dtos.PrepareChatSessionSuccessEventPayload, plugin_dtos.PrepareChatSessionFailedEventPayload](
		p,
		ctx,
		plugin_dtos.PrepareChatSessionEvent,
		payload,
		plugin_dtos.PrepareChatSessionFailedEvent,
		plugin_dtos.PrepareChatSessionSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to prepare chat session", "error", err)
		return "", err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to prepare chat session", "error", resp.Error)
		return "", errors.New(resp.Error.Cause)
	}

	return resp.Data.SessionID, nil
}

func (p *PluginClient) SendChatMessage(ctx context.Context, sessionID string, text string) error {
	payload := plugin_dtos.SendChatMessageEventPayload{
		SessionID: sessionID,
		Text:      text,
	}

	resp, err := sendAndWait[plugin_dtos.SendChatMessageSuccessEventPayload, plugin_dtos.SendChatMessageFailedEventPayload](
		p,
		ctx,
		plugin_dtos.SendChatMessageEvent,
		payload,
		plugin_dtos.SendChatMessageFailedEvent,
		plugin_dtos.SendChatMessageSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to send chat session messages", "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to send chat session messages", "error", resp.Error)
		return errors.New(resp.Error.Cause)
	}

	return nil
}

func (p *PluginClient) ResetChatChannel(ctx context.Context, channelID string) error {
	payload := plugin_dtos.ResetChatChannelEventPayload{
		ChannelID: channelID,
	}

	resp, err := sendAndWait[plugin_dtos.ResetChatChannelSuccessEventPayload, plugin_dtos.ResetChatChannelFailedEventPayload](
		p,
		ctx,
		plugin_dtos.ResetChatChannelEvent,
		payload,
		plugin_dtos.ResetChatChannelFailedEvent,
		plugin_dtos.ResetChatChannelSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to reset chat channel", "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to reset chat channel", "error", resp.Error)
		return errors.New(resp.Error.Cause)
	}

	return nil
}
