package sdk

import "errors"

var (
	ErrConnectionFailed         = errors.New("connection failed")
	ErrConnectionNotFound       = errors.New("connection not found")
	ErrConnectionNotReady       = errors.New("connection not ready")
	ErrCreatePayloadFailed      = errors.New("create payload failed")
	ErrHandshakeFailed          = errors.New("handshake failed")
	ErrHandshakeTimeout         = errors.New("handshake timeout")
	ErrHandshakeCancelled       = errors.New("handshake cancelled")
	ErrPrepareChatChannelFailed = errors.New("prepare chat channel failed")
	ErrIncorrectToolArgs        = errors.New("incorrect tool args")
	ErrMarshalToolResultFailed  = errors.New("marshal results failed")
	ErrSystemPromptTooLarge     = errors.New("system prompt too large")
	ErrTimeout                  = errors.New("timeout")
	ErrActionCancelled          = errors.New("action cancelled")
)
