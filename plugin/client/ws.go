package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/ws"
)

func (p *PluginClient) On(event dto.WsEvent, handlerID string, handler EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.handlers[event] == nil {
		p.handlers[event] = make(map[string]EventHandler)
	}
	p.handlers[event][handlerID] = handler
	p.Logger.Debug("Registered event listener", "event", event, "handler_id", handlerID)
}

func (p *PluginClient) Off(event dto.WsEvent, handlerID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if handlers, exists := p.handlers[event]; exists {
		delete(handlers, handlerID)
		p.Logger.Debug("Removed event listener", "event", event, "handler_id", handlerID)

		if len(handlers) == 0 {
			delete(p.handlers, event)
		}
	}
}

func (p *PluginClient) dispatchEvent(msg *dto.WsMessage[any]) {
	p.mu.RLock()
	var currentHandlers []EventHandler
	if eventMap, exists := p.handlers[msg.Event]; exists {
		currentHandlers = make([]EventHandler, 0, len(eventMap))
		for _, handler := range eventMap {
			currentHandlers = append(currentHandlers, handler)
		}
	}
	p.mu.RUnlock()

	for _, handler := range currentHandlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					p.Logger.Error("Plugin event handler panicked", "event", msg.Event, "error", r)
				}
			}()
			handler(msg)
		}()
	}
}

func WsRequest[TResponse any, TError any](
	p *PluginClient,
	sendEvent dto.WsEvent,
	successEvent dto.WsEvent,
	failureEvent dto.WsEvent,
	payload any,
	timeoutMs time.Duration,
) (*TResponse, *TError, error) {
	if timeoutMs == 0 {
		timeoutMs = 10 * time.Second
	}

	successHandlerID := uuid.New().String()
	failedHandlerID := uuid.New().String()

	respChan := make(chan *TResponse, 1)
	errChan := make(chan *TError, 1)

	var successHandler EventHandler
	successHandler = func(msg *dto.WsMessage[any]) {
		p.Off(successEvent, successHandlerID)
		p.Off(failureEvent, failedHandlerID)

		var target TResponse
		bytesData, err := json.Marshal(msg.Payload)
		if err != nil {
			p.Logger.Error("Failed to marshal ws response data", "error", err)
			return
		}
		if err := json.Unmarshal(bytesData, &target); err != nil {
			p.Logger.Error("Failed to unmarshal ws response data", "error", err)
			return
		}

		respChan <- &target
	}

	var failureHandler EventHandler
	failureHandler = func(msg *dto.WsMessage[any]) {
		p.Off(successEvent, successHandlerID)
		p.Off(failureEvent, failedHandlerID)

		var target TError
		bytesData, err := json.Marshal(msg.Payload)
		if err != nil {
			p.Logger.Error("Failed to marshal ws error data", "error", err)
			return
		}
		if err := json.Unmarshal(bytesData, &target); err != nil {
			p.Logger.Error("Failed to unmarshal ws error data", "error", err)
			return
		}

		errChan <- &target
	}

	p.On(successEvent, successHandlerID, successHandler)
	p.On(failureEvent, failedHandlerID, failureHandler)
	err := ws.SendMessage(p.ws, sendEvent, payload)
	if err != nil {
		p.Off(successEvent, successHandlerID)
		p.Off(failureEvent, failedHandlerID)
		return nil, nil, fmt.Errorf("failed to send ws request: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutMs)
	defer cancel()

	select {
	case <-ctx.Done():
		p.Off(successEvent, successHandlerID)
		p.Off(failureEvent, failedHandlerID)
		return nil, nil, fmt.Errorf("websocket request timeout for event: %v", sendEvent)

	case res := <-respChan:
		return res, nil, nil

	case terr := <-errChan:
		return nil, terr, nil
	}
}
