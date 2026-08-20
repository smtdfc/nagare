package handlers

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/messages"
	"github.com/smtdfc/nagare/shared/ws"
)

type ChatHandler struct{}

var (
	pendingMutex     sync.Mutex
	pendingResponses = make(map[string]bool)
)

func (h *ChatHandler) CreateSession(i *ws.WsInstance) {
	if i.Auth == nil {
		ws.SendMessage(i, dto.WS_CREATE_SESSION_FAILED, dto.CreateSessionFailedEvent{
			Cause: "Unauthorized",
		})
		return
	}

	sessionID, err := global.GlobalSessionMgr.CreateSession()
	if err != nil {
		ws.SendMessage(i, dto.WS_CREATE_SESSION_FAILED, dto.CreateSessionFailedEvent{
			Cause: err.Error(),
		})
		return
	}

	ws.SendMessage(i, dto.WS_CREATE_SESSION_SUCCESS, dto.CreateSessionSuccessEvent{
		ID: sessionID,
	})
}

func (h *ChatHandler) InvokeAgent(i *ws.WsInstance, message *dto.WsMessage[any]) {
	if i.Auth == nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    "",
			Cause: "Unauthorized",
		})
		return
	}

	pendingMutex.Lock()
	if pendingResponses[i.ID] {
		pendingMutex.Unlock()
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    "",
			Cause: "Request rejected: The agent is busy processing an active response. Please wait until it completes.",
		})
		return
	}
	pendingResponses[i.ID] = true
	pendingMutex.Unlock()

	defer func() {
		pendingMutex.Lock()
		pendingResponses[i.ID] = false
		pendingMutex.Unlock()
	}()

	payload, err := ws.GetPayload[dto.InvokeAgentEvent](message)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    "",
			Cause: "Invalid payload",
		})

		return
	}

	id := payload.ID
	sessionID := payload.SessionID
	text := payload.Text

	history, err := global.GlobalSessionMgr.GetMessagesBySessionID(sessionID)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    id,
			Cause: fmt.Sprintf("Failed to fetch messages for session ID: %s", sessionID),
		})

		return
	}

	state := global.CreateEmptyAgentState().WithHistory(history)
	agent, err := global.FetchReadyAgent(state)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    id,
			Cause: err.Error(),
		})

		return
	}

	pendingResponses[i.ID] = true
	output, err := agent.Invoke(messages.NewText(text, messages.USER))
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    id,
			Cause: "Execution failed: An unexpected error occurred during the agent's ReAct reasoning loop",
		})

		return
	}

	for chunk := range output {
		ws.SendMessage(i, dto.WS_AGENT_RESPONSE, dto.AgentOutputEvent{
			ID:        id,
			SessionID: sessionID,
			Message:   chunk,
		})
	}

	err = global.GlobalSessionMgr.SaveSession(sessionID, state.PendingMessages)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    id,
			Cause: "Failed to save history",
		})

		return
	}

	err = state.CommitMessage()
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailedEvent{
			ID:    id,
			Cause: "Failed to save history",
		})

		return
	}

	global.PutAgentIntoPool(agent)
}

// @Injectable
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}
