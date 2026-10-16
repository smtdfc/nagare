package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

func GetData[T any](payload *plugin_dtos.Payload[any]) (*T, error) {
	if payload == nil || payload.Data == nil {
		return nil, fmt.Errorf("payload or data is nil")
	}

	bytesData, err := json.Marshal(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload data: %w", err)
	}

	var result T
	if err := json.Unmarshal(bytesData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into target type %T: %w", (*T)(nil), err)
	}

	return &result, nil
}

type plugin_dtosActionResult[T any, F any] struct {
	IsSuccess bool
	Data      *T
	Error     *F
}

func sendAndWait[S any, F any, P any](
	p *PluginClient,
	ctx context.Context,
	event plugin_dtos.Event,
	payload P,
	failedEvent plugin_dtos.Event,
	successEvent plugin_dtos.Event,
) (*plugin_dtosActionResult[S, F], error) {
	requestID := uuid.New().String()
	respChan := make(chan *plugin_dtos.Payload[any], 1)

	p.mu.Lock()
	p.pendingRequests[requestID] = respChan
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.pendingRequests, requestID)
		p.mu.Unlock()
	}()

	err := p.connector.Send(event, payload, requestID)
	if err != nil {
		p.Logger.Error("failed to send event", "event", event, "error", err)
		return nil, err
	}

	result := &plugin_dtosActionResult[S, F]{
		IsSuccess: false,
		Data:      nil,
		Error:     nil,
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-respChan:
		if resp.Event == failedEvent {
			result.IsSuccess = false
			failPayload, _ := GetData[F](resp)
			result.Error = failPayload
			return result, nil
		}
		if successEvent != "" && resp.Event == successEvent {
			result.IsSuccess = true
			successPayload, _ := GetData[S](resp)
			result.Data = successPayload
			return result, nil
		}
		p.Logger.Error("unexpected event received", "event", resp.Event)
		return nil, errors.New("unexpected event received")
	}
}
