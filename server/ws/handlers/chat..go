package handlers

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/server/ws"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/messages"
)

type ChatHandler struct{}

var (
	pendingMutex     sync.Mutex
	pendingResponses = make(map[string]bool)
)

func (h *ChatHandler) CreateSession(i *ws.WsInstance) {
	sessionID, err := global.GlobalSessionMgr.CreateSession()
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
	pendingMutex.Lock()
	if pendingResponses[i.ID] {
		pendingMutex.Unlock()
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
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

	payload, err := ws.GetPayload[dto.InvokeAgent](message)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
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
		fmt.Println(err.Error())
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: fmt.Sprintf("Failed to fetch messages for session ID: %s", sessionID),
		})

		return
	}

	state := global.CreateEmptyAgentState().WithHistory(history)
	agent, err := global.FetchReadyAgent(state)
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: err.Error(),
		})

		return
	}

	pendingResponses[i.ID] = true
	output, err := agent.Invoke(messages.NewText(text, messages.USER))
	if err != nil {
		ws.SendMessage(i, dto.WS_INVOKE_AGENT_FAILED, dto.InvokeAgentFailed{
			ID:    id,
			Cause: "Execution failed: An unexpected error occurred during the agent's ReAct reasoning loop",
		})

		return
	}

	for chunk := range output {
		ws.SendMessage(i, dto.WS_AGENT_RESPONSE, dto.AgentOuput{
			ID:        id,
			SessionID: sessionID,
			Message:   chunk,
		})
	}

	err = global.GlobalSessionMgr.SaveSession(sessionID, state.PendingMessages)
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

	global.PutAgentIntoPool(agent)
}

// @Injectable
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}
