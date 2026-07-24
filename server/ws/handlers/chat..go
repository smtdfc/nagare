package handlers

import (
	"fmt"

	"github.com/smtdfc/nagare/core/providers"
	"github.com/smtdfc/nagare/server/ws"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/messages"
)

type ChatHandler struct{}

func (h *ChatHandler) CreateSession(i *ws.WsInstance) {
	sessionID, err := providers.GlobalSessionManager.CreateSession()
	if err != nil {
		ws.SendMessage(i, dto.WS_CREATE_SESSION_FAILED, dto.CreateSessionFailed{
			Cause: err.Error(),
		})
		return
	}

	ws.SendMessage(i, dto.WS_CREATE_SESSION_SUCCESS, dto.CreateSessionSuccess{
		ID: sessionID,
	})
}

func (h *ChatHandler) InvokeAgent(i *ws.WsInstance, message *dto.WsMessage[any]) {
	payload, err := ws.GetPayload[dto.InvokeAgent](message)
	if err != nil {
		fmt.Println(err.Error())
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    "",
			Cause: "Invalid payload",
		})

		return
	}

	id := payload.ID
	sessionID := payload.SessionID
	text := payload.Text

	history, err := providers.GlobalSessionManager.GetMessagesBySessionID(sessionID)
	if err != nil {
		fmt.Println(err.Error())
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: fmt.Sprintf("Failed to fetch messages for session ID: %s", sessionID),
		})

		return
	}

	state := providers.CreateEmptyAgentState().WithHistory(history)
	agent, err := providers.FetchReadyAgent(state)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: "The agent is currently initializing or busy and cannot process requests",
		})

		return
	}

	output, err := agent.Invoke(messages.NewText(text, messages.USER))
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: "Execution failed: An unexpected error occurred during the agent's ReAct reasoning loop",
		})

		return
	}

	for chunk := range output {
		ws.SendMessage(i, dto.WS_AGENT_RESPONSE, dto.AgentResponse{
			ID:        id,
			SessionID: sessionID,
			Message:   chunk,
		})
	}

	err = providers.GlobalSessionManager.SaveSession(sessionID, state.PendingMessages)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: "Failed to save history",
		})

		return
	}

	err = state.CommitMessage()
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: "Failed to save history",
		})

		return
	}
}

// @Injectable
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}
