package chat

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/chat"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/shared/helpers"

	websocket_dtos "github.com/smtdfc/nagare/shared/dtos/websocket"
)

type Worker struct {
	mu           sync.RWMutex
	chatEventBus *EventBus
	ws           *websocket.Coordinator
	sessionMgr   *manager.SessionManager
	agentInvoker *chat.AgentInvoker
}

func (c *Worker) HandleChat(payload *SendMessageEvent) {
	output, _ := c.agentInvoker.Invoke(
		payload.SessionID,
		payload.Text,
		payload.SenderType,
		payload.SenderID,
	)

	for chunk := range output {
		chunkJson, _ := helpers.MarshalJson(chunk)
		err := websocket.BroadcastToRoom(
			c.ws,
			fmt.Sprintf("session:%s", payload.SessionID),
			websocket_dtos.ReceivedChatMessageEvent,
			&websocket_dtos.ReceivedChatMessageEventPayload{
				SessionID: payload.SessionID,
				Message:   chunkJson,
			},
			payload.RequestID,
			nil,
		)
		if err != nil {
			return
		}
	}
}

func (c *Worker) Handle(evt EventPayload) {
	switch evt.GetEventType() {
	case SendEvent:
		payload := evt.(*SendMessageEvent)
		c.HandleChat(payload)
	}
}

func (c *Worker) Start() {
	go func() {
		ch, unsubscribe := c.chatEventBus.Subscribe(string(Topic))
		defer unsubscribe()
		for chatEventPayload := range ch {
			go c.Handle(chatEventPayload)
		}
	}()
}

// @Injectable
func NewWorker(chatEventBus *EventBus, ws *websocket.Coordinator, sessionMgr *manager.SessionManager, agentInvoker *chat.AgentInvoker) *Worker {
	return &Worker{
		chatEventBus: chatEventBus,
		mu:           sync.RWMutex{},
		ws:           ws,
		sessionMgr:   sessionMgr,
		agentInvoker: agentInvoker,
	}
}
