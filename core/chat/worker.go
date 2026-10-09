package chat

import (
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
)

type Worker struct {
	chatEventBus *event_bus.CoreEventBus
	agentInvoker *AgentInvoker
	logger       *logger.BaseLogger
}

func (c *Worker) Handle(evt event_bus.EventPayload) {
	switch evt.GetEventType() {
	case event_bus.SendEvent:
		payload := evt.(*event_bus.SendMessageEventPayload)
		c.HandleSendMessageEvent(payload)
	}
}

func (c *Worker) Do() {
	go func() {
		ch, unsubscribe := c.chatEventBus.Subscribe(event_bus.SendEvent)
		defer unsubscribe()
		for chunkEventPayload := range ch {
			go c.Handle(chunkEventPayload)
		}
	}()
}

func (c *Worker) HandleSendMessageEvent(payload *event_bus.SendMessageEventPayload) {
	output, err := c.agentInvoker.Invoke(&AgentInvokeParams{
		SessionID:        payload.SessionID,
		InputMessages:    payload.Messages,
		SenderType:       payload.SenderType,
		SenderID:         payload.SenderID,
		SendIntoEventBus: true,
		InvokeID:         payload.InvokeID,
	})
	if err != nil {
		return
	}

	go func() {
		for range output {
		}
	}()
}

// @Injectable
func NewChatWorker(eventBus *event_bus.CoreEventBus, agentInvoker *AgentInvoker, logger *logger.BaseLogger) *Worker {
	return &Worker{
		chatEventBus: eventBus,
		agentInvoker: agentInvoker,
		logger:       logger.With("worker", "core:chat"),
	}
}
