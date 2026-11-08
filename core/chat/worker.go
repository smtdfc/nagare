package chat

import (
	"context"
	"sync"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
)

type Worker struct {
	chatEventBus *event_bus.CoreEventBus
	agentInvoker *AgentInvoker
	logger       *logger.BaseLogger
	cancelFnMu   sync.Mutex
	cancelFn     map[string]context.CancelFunc
}

func (c *Worker) Handle(evt event_bus.EventPayload) {
	switch evt.GetEventType() {
	case event_bus.SendChatMessageEvent:
		payload := evt.(*event_bus.SendChatMessageEventPayload)
		c.HandleSendMessageEvent(payload)
	case event_bus.CancelChatEvent:
		payload := evt.(*event_bus.CancelChatEventPayload)
		c.HandleCancelEvent(payload)
	case event_bus.CompactConversationEvent:
		payload := evt.(*event_bus.CompactConversationEventPayload)
		c.HandleCompactConversationEvent(payload)
	}
}

func (c *Worker) subscribeAndHandle(event event_bus.EventType) {
	go func() {
		ch, unsubscribe := c.chatEventBus.Subscribe(event)
		defer unsubscribe()

		for payload := range ch {
			go c.Handle(payload)
		}
	}()
}

func (c *Worker) Do() {
	c.subscribeAndHandle(event_bus.SendChatMessageEvent)
	c.subscribeAndHandle(event_bus.CancelChatEvent)
	c.subscribeAndHandle(event_bus.CompactConversationEvent)
}

func (c *Worker) HandleSendMessageEvent(payload *event_bus.SendChatMessageEventPayload) {
	ctx, cancelFn := context.WithCancel(context.Background())
	c.cancelFnMu.Lock()
	c.cancelFn[payload.InvokeID] = cancelFn
	c.cancelFnMu.Unlock()

	output, err := c.agentInvoker.Invoke(ctx, &AgentInvokeParams{
		SessionID:        payload.SessionID,
		InputMessages:    payload.Messages,
		SenderType:       payload.SenderType,
		SenderID:         payload.SenderID,
		SendIntoEventBus: true,
		InvokeID:         payload.InvokeID,
		ContextIsolate:   payload.ContextIsolate,
	})
	if err != nil {
		c.cancelFnMu.Lock()
		delete(c.cancelFn, payload.InvokeID)
		c.cancelFnMu.Unlock()
		return
	}

	go func() {
		defer func() {
			c.cancelFnMu.Lock()
			delete(c.cancelFn, payload.InvokeID)
			c.cancelFnMu.Unlock()
		}()
		for range output {
		}
	}()
}

func (c *Worker) HandleCancelEvent(payload *event_bus.CancelChatEventPayload) {
	c.cancelFnMu.Lock()
	cancelFn, ok := c.cancelFn[payload.InvokeID]
	if ok {
		delete(c.cancelFn, payload.InvokeID)
	}
	c.cancelFnMu.Unlock()

	if ok && cancelFn != nil {
		cancelFn()
	}
}

func (c *Worker) HandleCompactConversationEvent(payload *event_bus.CompactConversationEventPayload) {
	ctx, cancelFn := context.WithCancel(context.Background())
	c.cancelFnMu.Lock()
	c.cancelFn[payload.InvokeID] = cancelFn
	c.cancelFnMu.Unlock()

	output, err := c.agentInvoker.Compact(ctx, &ConversationCompactParams{
		SessionID:        payload.SessionID,
		SenderType:       payload.SenderType,
		SenderID:         payload.SenderID,
		SendIntoEventBus: true,
		InvokeID:         payload.InvokeID,
	})
	if err != nil {
		c.cancelFnMu.Lock()
		delete(c.cancelFn, payload.InvokeID)
		c.cancelFnMu.Unlock()
		return
	}

	go func() {
		defer func() {
			c.cancelFnMu.Lock()
			delete(c.cancelFn, payload.InvokeID)
			c.cancelFnMu.Unlock()
		}()
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
		cancelFn:     make(map[string]context.CancelFunc),
	}
}
