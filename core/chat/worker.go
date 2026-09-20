package chat

import (
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
)

type ChatWorker struct {
	chatEventBus *event_bus.CoreEventBus
	agentInvoker *AgentInvoker
	logger       *logger.BaseLogger
}

func (c *ChatWorker) Handle(evt event_bus.EventPayload) {
	switch evt.GetEventType() {
	case event_bus.SendEvent:
		payload := evt.(*event_bus.SendMessageEventPayload)
		c.HandleSendMessageEvent(payload)
	}
}

func (w *ChatWorker) Do() {
	go func() {
		ch, unsubscribe := w.chatEventBus.Subscribe(event_bus.SendEvent)
		defer unsubscribe()
		for chunkEventPayload := range ch {
			go w.Handle(chunkEventPayload)
		}
	}()
}

func (w *ChatWorker) HandleSendMessageEvent(payload *event_bus.SendMessageEventPayload) {
	output, err := w.agentInvoker.Invoke(
		payload.SessionID,
		payload.Text,
		payload.SenderType,
		payload.SenderID,
		true,
	)
	if err != nil {
		return
	}

	go func() {
		for _ = range output {
		}
	}()
}

// @Injectable
func NewChatWorker(eventBus *event_bus.CoreEventBus, agentInvoker *AgentInvoker, logger *logger.BaseLogger) *ChatWorker {
	return &ChatWorker{
		chatEventBus: eventBus,
		agentInvoker: agentInvoker,
		logger:       logger.With("worker", "core:chat"),
	}
}
