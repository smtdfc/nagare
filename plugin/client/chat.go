package client

import (
	"context"
	"errors"
	"uuid"

	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
)

func (p *PluginClient) PrepareChatSession(ctx context.Context, channelID string) (string, error) {
	requestID := uuid.NewV4().String()
	respChan := make(chan *websocket.Payload[any], 1)

	p.mu.Lock()
	p.pendingRequests[requestID] = respChan
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.pendingRequests, requestID)
		p.mu.Unlock()
	}()

	err := p.connector.Send(
		plugin_dtos.PrepareChatSessionEvent,
		plugin_dtos.PrepareChatSessionEventPayload{
			ChannelID: channelID,
		},
		requestID,
	)
	if err != nil {
		p.Logger.Error("failed to send prepare chat session event", "error", err)
		return "", err
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case resp := <-respChan:
		if resp.Event == plugin_dtos.PrepareChatSessionFailedEvent {
			payload, _ := GetData[plugin_dtos.PrepareChatSessionFailedEventPayload](resp)
			p.Logger.Error("failed to prepare chat session event", "error", payload)
			return "", errors.New(payload.Cause)
		}

		if resp.Event == plugin_dtos.PrepareChatSessionSuccessEvent {
			payload, _ := GetData[plugin_dtos.PrepareChatSessionSuccessEventPayload](resp)
			return payload.SessionID, nil
		}

		return "", nil
	}

}

func (p *PluginClient) SendChatMessage(ctx context.Context, sessionID string, text string) error {
	requestID := uuid.NewV4().String()
	respChan := make(chan *websocket.Payload[any], 1)

	p.mu.Lock()
	p.pendingRequests[requestID] = respChan
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.pendingRequests, requestID)
		p.mu.Unlock()
	}()

	err := p.connector.Send(
		plugin_dtos.SendChatMessageEvent,
		plugin_dtos.SendChatMessageEventPayload{
			SessionID: sessionID,
			Text:      text,
		},
		requestID,
	)
	if err != nil {
		p.Logger.Error("failed to send chat session message", "error", err)
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-respChan:
		if resp.Event == plugin_dtos.SendChatMessageFailedEvent {
			payload, _ := GetData[plugin_dtos.SendChatMessageFailedEventPayload](resp)
			p.Logger.Error("failed to send chat session message", "error", payload)
			return errors.New(payload.Cause)
		}

		return nil
	}

}
