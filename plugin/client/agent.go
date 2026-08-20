package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/messages"
	"github.com/smtdfc/nagare/shared/ws"
)

func (p *PluginClient) InvokeAgent(text string, sessionID string) (<-chan messages.Message, error) {
	output := make(chan messages.Message)
	requestID := uuid.New().String()
	successHandlerID := uuid.New().String()
	failedHandlerID := uuid.New().String()

	var once sync.Once
	var safeClose func()

	handleInvokeAgentFailed := func(msg *dto.WsMessage[any]) {
		payload, err := ws.GetPayload[dto.PluginInvokeAgentFailedEvent](msg)
		if err != nil {
			p.Logger.Error("Failed to get payload", "error", err, "session_id", sessionID)
			safeClose()
			return
		}

		if requestID != payload.ID {
			return
		}

		output <- messages.NewResponseFailed("400", payload.Cause)
		safeClose()
	}

	handleInvokeAgentSuccess := func(msg *dto.WsMessage[any]) {
		payload, err := ws.GetPayload[dto.PluginInvokeAgentSuccessEvent](msg)
		if err != nil {
			p.Logger.Error("Failed to get payload", "error", err, "session_id", sessionID)
			safeClose()
			return
		}

		if requestID != payload.ID {
			return
		}

		messageBytes, err := json.Marshal(payload.Message)
		if err != nil {
			return
		}

		message, err := helpers.ParseAgentMessage(messageBytes)
		if err != nil {
			p.Logger.Error("Failed to parse message chunk", "error", err, "session_id", sessionID)
			safeClose()
			return
		}
		output <- message

		switch chunk := message.(type) {
		case *messages.AgentResponse:
			if chunk.Status == messages.AGENT_RESPONSE_COMPLETED || chunk.Status == messages.AGENT_RESPONSE_FAILED {
				safeClose()
			}
		}
	}

	safeClose = func() {
		once.Do(func() {
			close(output)
			p.Off(dto.WS_PLUGIN_INVOKE_AGENT_SUCCESS, successHandlerID)
			p.Off(dto.WS_PLUGIN_INVOKE_AGENT_FAILED, failedHandlerID)
		})
	}

	go (func() {
		p.On(dto.WS_PLUGIN_INVOKE_AGENT_SUCCESS, successHandlerID, handleInvokeAgentSuccess)
		p.On(dto.WS_PLUGIN_INVOKE_AGENT_FAILED, failedHandlerID, handleInvokeAgentFailed)
	})()

	err := ws.SendMessage(
		p.ws,
		dto.WS_PLUGIN_INVOKE_AGENT,
		dto.PluginInvokeAgentEvent{
			ID:        requestID,
			SessionID: sessionID,
			Text:      text,
		},
	)

	fmt.Println(err)
	return output, nil
}

func (p *PluginClient) ResetChatSession(sessionID string) error {
	_, werr, err := WsRequest[dto.PluginResetChatSessionSuccessEvent, dto.PluginResetChatSessionFailedEvent](
		p,
		dto.WS_PLUGIN_RESET_CHAT_SESSION,
		dto.WS_PLUGIN_RESET_CHAT_SESSION_SUCCESS,
		dto.WS_PLUGIN_RESET_CHAT_SESSION_FAILED,
		&dto.PluginResetChatSessionEvent{
			ID:        sessionID,
			SessionID: sessionID,
		},
		10*time.Minute,
	)

	if err != nil {
		return err
	}

	if werr != nil {
		return errors.New(werr.Cause)
	}

	return nil
}
