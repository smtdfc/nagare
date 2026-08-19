package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/ws"
)

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

	respChan := make(chan *TResponse, 1)
	errChan := make(chan *TError, 1)

	var successHandler EventHandler
	successHandler = func(msg *dto.WsMessage[any]) {
		p.Off(successEvent)
		p.Off(failureEvent)

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
		p.Off(successEvent)
		p.Off(failureEvent)

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

	p.On(successEvent, successHandler)
	p.On(failureEvent, failureHandler)
	err := ws.SendMessage(p.ws, sendEvent, payload)
	if err != nil {
		p.Off(successEvent)
		p.Off(failureEvent)
		return nil, nil, fmt.Errorf("failed to send ws request: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutMs)
	defer cancel()

	select {
	case <-ctx.Done():
		p.Off(successEvent)
		p.Off(failureEvent)
		return nil, nil, fmt.Errorf("websocket request timeout for event: %v", sendEvent)

	case res := <-respChan:
		return res, nil, nil

	case terr := <-errChan:
		return nil, terr, nil
	}
}
