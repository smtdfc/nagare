package client

import (
	"context"
	"errors"

	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

func (p *PluginClient) PrepareChatChannel(ctx context.Context, channelID string) error {
	payload := plugin_dtos.PrepareChatChannelEventPayload{
		ChannelID: channelID,
	}

	resp, err := sendAndWait[plugin_dtos.PrepareChatChannelSuccessEventPayload, plugin_dtos.PrepareChatChannelFailedEventPayload](
		p,
		ctx,
		plugin_dtos.PrepareChatChannelEvent,
		payload,
		plugin_dtos.PrepareChatChannelFailedEvent,
		plugin_dtos.PrepareChatChannelSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to prepare chat channel", "channelID", channelID, "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to prepare chat channel", "channelID", channelID, "error", resp.Error)
		return errors.New(resp.Error.Cause)
	}

	return nil
}

func (p *PluginClient) SendChatMessage(ctx context.Context, channelID string, text string) error {
	payload := plugin_dtos.SendChatMessageEventPayload{
		ChannelID: channelID,
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
		p.Logger.Error("failed to send chat message", "channelID", channelID, "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to send chat message", "channelID", channelID, "error", resp.Error)
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

func (p *PluginClient) SubmitToolApproval(ctx context.Context, channelID, approvalRequestID string, accept bool) error {
	resp, err := sendAndWait[plugin_dtos.SubmitToolApprovalSuccessEventPayload, plugin_dtos.SubmitToolApprovalFailedEventPayload](
		p,
		ctx,
		plugin_dtos.SubmitToolApprovalEvent,
		plugin_dtos.SubmitToolApprovalEventPayload{
			ChannelID:         channelID,
			ApprovalRequestID: approvalRequestID,
			Accept:            accept,
		},
		plugin_dtos.SubmitToolApprovalFailedEvent,
		plugin_dtos.SubmitToolApprovalSuccessEvent,
	)
	if err != nil {
		return err
	}
	if !resp.IsSuccess {
		return errors.New(resp.Error.Cause)
	}
	return nil
}
