package external_module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	InvokePrefix string = "module:invoke"
	ReplyPrefix  string = "module:reply"
)

func MakeInvokeEvent(action string) Event {
	if action == "" {
		return Event(InvokePrefix)
	}
	return Event(fmt.Sprintf("%s:%s", InvokePrefix, action))
}

func MakeReplyEvent(action string) Event {
	if action == "" {
		return Event(ReplyPrefix)
	}
	return Event(fmt.Sprintf("%s:%s", ReplyPrefix, action))
}

func ParseActionEvent(event Event) (prefix string, action string) {
	str := string(event)
	if strings.HasPrefix(str, InvokePrefix) {
		trimmed := strings.TrimPrefix(str, InvokePrefix)
		trimmed = strings.TrimPrefix(trimmed, ":")
		return InvokePrefix, trimmed
	}
	if strings.HasPrefix(str, ReplyPrefix) {
		trimmed := strings.TrimPrefix(str, ReplyPrefix)
		trimmed = strings.TrimPrefix(trimmed, ":")
		return ReplyPrefix, trimmed
	}
	return "", ""
}

type InvokePayload[T any] struct {
	Action    string `json:"action,omitempty"`
	RequestID string `json:"requestID,omitempty"`
	Data      T      `json:"data"`
}

func (p *InvokePayload[T]) BindData(target any) error {
	if any(p.Data) == nil {
		return errors.New("data is nil")
	}

	bytesData, err := json.Marshal(p.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal invoke payload data: %w", err)
	}

	if err := json.Unmarshal(bytesData, target); err != nil {
		return fmt.Errorf("failed to unmarshal into target %T: %w", target, err)
	}

	return nil
}

type ReplyPayload[T any] struct {
	Action    string `json:"action,omitempty"`
	RequestID string `json:"requestID,omitempty"`
	Data      T      `json:"data"`
}

func (p *ReplyPayload[T]) BindData(target any) error {
	if any(p.Data) == nil {
		return errors.New("data is nil")
	}

	bytesData, err := json.Marshal(p.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal reply payload data: %w", err)
	}

	if err := json.Unmarshal(bytesData, target); err != nil {
		return fmt.Errorf("failed to unmarshal into target %T: %w", target, err)
	}

	return nil
}

func DecodePayloadData[U any, T any](data T) (*U, error) {
	if any(data) == nil {
		return nil, errors.New("data is nil")
	}

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var target U
	if err := json.Unmarshal(bytesData, &target); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into type %T: %w", target, err)
	}

	return &target, nil
}

type ModuleChannelPair[T any, R any] struct {
	InvokeChan chan<- InvokePayload[T]
	ReplyChan  <-chan ReplyPayload[R]
}

func (h *ExternalModuleHost) RegisterChannels(moduleID string, bufferSize int) (chan<- InvokePayload[any], <-chan ReplyPayload[any], error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if bufferSize <= 0 {
		bufferSize = 64
	}

	if _, exists := h.invokeChannels[moduleID]; exists {
		return nil, nil, fmt.Errorf("channels for module %s already registered", moduleID)
	}

	invokeCh := make(chan InvokePayload[any], bufferSize)
	replyCh := make(chan ReplyPayload[any], bufferSize)

	h.invokeChannels[moduleID] = invokeCh
	h.replyChannels[moduleID] = replyCh

	go h.forwardToModule(moduleID, invokeCh)

	return invokeCh, replyCh, nil
}

func (h *ExternalModuleHost) UnregisterChannels(moduleID string) {
	h.mu.Lock()
	invokeCh, hasInvoke := h.invokeChannels[moduleID]
	delete(h.invokeChannels, moduleID)

	replyCh, hasReply := h.replyChannels[moduleID]
	delete(h.replyChannels, moduleID)
	h.mu.Unlock()

	if hasInvoke && invokeCh != nil {
		close(invokeCh)
	}
	if hasReply && replyCh != nil {
		close(replyCh)
	}
}

func (h *ExternalModuleHost) forwardToInternal(pConn *Connection, payload *Payload[any]) {
	h.mu.RLock()
	channel, exist := h.replyChannels[pConn.moduleID]
	h.mu.RUnlock()

	if !exist {
		h.logger.Warn("No reply consumer channel registered for module", "moduleID", pConn.moduleID)
		return
	}

	_, action := ParseActionEvent(payload.Event)

	select {
	case channel <- ReplyPayload[any]{
		Action:    action,
		RequestID: payload.RequestID,
		Data:      payload.Data,
	}:
	default:
		h.logger.Warn("Reply channel buffer full, dropping message", "moduleID", pConn.moduleID, "requestID", payload.RequestID)
	}
}

func (h *ExternalModuleHost) forwardToModule(moduleID string, invokeCh <-chan InvokePayload[any]) {
	for invokePayload := range invokeCh {
		h.mu.RLock()
		conn, ok := h.connections[moduleID]
		h.mu.RUnlock()

		if !ok || conn == nil {
			h.logger.Warn("External module is not connected, dropping invoke payload", "moduleID", moduleID)
			continue
		}

		event := MakeInvokeEvent(invokePayload.Action)
		if err := conn.Send(event, invokePayload.Data, invokePayload.RequestID); err != nil {
			h.logger.Error("Failed to send invoke payload to external module", "moduleID", moduleID, "error", err)
		}
	}
}

func RequestReply[T any, R any](
	ctx context.Context,
	invokeCh chan<- InvokePayload[any],
	replyCh <-chan ReplyPayload[any],
	action string,
	data T,
	requestID string,
) (res *ReplyPayload[R], err error) {
	if invokeCh == nil {
		return nil, errors.New("invoke channel is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invoke channel send panic: %v", r)
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case invokeCh <- InvokePayload[any]{
		Action:    action,
		RequestID: requestID,
		Data:      data,
	}:
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case reply, ok := <-replyCh:
			if !ok {
				return nil, errors.New("reply channel closed")
			}
			if (requestID == "" || reply.RequestID == requestID) && (action == "" || reply.Action == "" || reply.Action == action) {
				var target R
				if err := reply.BindData(&target); err != nil {
					return nil, fmt.Errorf("failed to bind reply data to target type %T: %w", target, err)
				}
				return &ReplyPayload[R]{
					Action:    reply.Action,
					RequestID: reply.RequestID,
					Data:      target,
				}, nil
			}
		}
	}
}
