package workers

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/chat"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/gateway/event"
	"github.com/smtdfc/nagare/gateway/websocket"
	"github.com/smtdfc/nagare/shared/event_bus"
	"github.com/smtdfc/nagare/shared/helpers"

	websocket_dtos "github.com/smtdfc/nagare/shared/dtos/websocket"
)

type ChatWorker struct {
	mu           sync.RWMutex
	bus          *event_bus.EventBus[event.ChatEventPayload]
	ws           *websocket.Coordinator
	sessionMgr   *manager.SessionManager
	agentInvoker *chat.AgentInvoker
}

func (c *ChatWorker) HandleChat(payload *event.ChatSendMessageEvent) {
	output, _ := c.agentInvoker.Invoke(
		payload.SessionID,
		payload.Text,
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
			nil,
		)
		if err != nil {
			return
		}
	}
}

func (c *ChatWorker) Handle(evt event.ChatEventPayload) {
	switch evt.GetEventType() {
	case event.ChatSendEvent:
		payload := evt.(*event.ChatSendMessageEvent)
		c.HandleChat(payload)
	}
}

func (c *ChatWorker) Start() {
	go func() {
		ch, unsubscribe := c.bus.Subscribe(string(event.ChatTopic))
		defer unsubscribe()
		for chatEventPayload := range ch {
			go c.Handle(chatEventPayload)
		}
	}()
}

// @Injectable
func NewChatWorker(busSys *event.AppEventBusSystem, ws *websocket.Coordinator, sessionMgr *manager.SessionManager, agentInvoker *chat.AgentInvoker) *ChatWorker {
	return &ChatWorker{
		bus:          busSys.ChatEventBus,
		mu:           sync.RWMutex{},
		ws:           ws,
		sessionMgr:   sessionMgr,
		agentInvoker: agentInvoker,
	}
}
