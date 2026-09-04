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
	ws           *websocket.WebsocketCoordinator
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
		websocket.BroadcastToRoom(
			c.ws,
			fmt.Sprintf("session:%s", payload.SessionID),
			websocket_dtos.CHAT_RECEIVED_MESSAGE_EVENT,
			&websocket_dtos.ChatReceivedMessageEvent{
				SessionID: payload.SessionID,
				Message:   chunkJson,
			},
			nil,
		)
	}
}

func (c *ChatWorker) Handle(evt event.ChatEventPayload) {
	switch evt.GetEventType() {
	case event.CHAT_SEND_EVENT:
		payload := evt.(*event.ChatSendMessageEvent)
		c.HandleChat(payload)
	}
}

func (c *ChatWorker) Start() {
	go func() {
		ch, unsubscribe := c.bus.Subscribe(string(event.CHAT_TOPIC))
		defer unsubscribe()
		for event := range ch {
			go c.Handle(event)
		}
	}()
}

// @Injectable
func NewChatWorker(busSys *event.AppEventBusSystem, ws *websocket.WebsocketCoordinator, sessionMgr *manager.SessionManager, agentInvoker *chat.AgentInvoker) *ChatWorker {
	return &ChatWorker{
		bus:          busSys.ChatEventBus,
		mu:           sync.RWMutex{},
		ws:           ws,
		sessionMgr:   sessionMgr,
		agentInvoker: agentInvoker,
	}
}
