package sdk

import (
	"context"
	"errors"

	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (p *PluginClient) PrepareChatChannel(ctx context.Context, channelID string) error {
	payload := contracts.PrepareChatChannelEventPayload{
		ChannelID: channelID,
	}

	resp, err := sendAndWait[contracts.PrepareChatChannelSuccessEventPayload, contracts.PrepareChatChannelFailedEventPayload](
		p,
		ctx,
		contracts.PrepareChatChannelEvent,
		payload,
		contracts.PrepareChatChannelFailedEvent,
		contracts.PrepareChatChannelSuccessEvent,
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
	payload := contracts.SendChatMessageEventPayload{
		ChannelID: channelID,
		Text:      text,
	}

	resp, err := sendAndWait[contracts.SendChatMessageSuccessEventPayload, contracts.SendChatMessageFailedEventPayload](
		p,
		ctx,
		contracts.SendChatMessageEvent,
		payload,
		contracts.SendChatMessageFailedEvent,
		contracts.SendChatMessageSuccessEvent,
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
	payload := contracts.ResetChatChannelEventPayload{
		ChannelID: channelID,
	}

	resp, err := sendAndWait[contracts.ResetChatChannelSuccessEventPayload, contracts.ResetChatChannelFailedEventPayload](
		p,
		ctx,
		contracts.ResetChatChannelEvent,
		payload,
		contracts.ResetChatChannelFailedEvent,
		contracts.ResetChatChannelSuccessEvent,
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
	resp, err := sendAndWait[contracts.SubmitToolApprovalSuccessEventPayload, contracts.SubmitToolApprovalFailedEventPayload](
		p,
		ctx,
		contracts.SubmitToolApprovalEvent,
		contracts.SubmitToolApprovalEventPayload{
			ChannelID:         channelID,
			ApprovalRequestID: approvalRequestID,
			Accept:            accept,
		},
		contracts.SubmitToolApprovalFailedEvent,
		contracts.SubmitToolApprovalSuccessEvent,
	)
	if err != nil {
		return err
	}
	if !resp.IsSuccess {
		return errors.New(resp.Error.Cause)
	}
	return nil
}
