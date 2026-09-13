package client

import "errors"

var (
	ErrConnectionFailed         = errors.New("connection failed")
	ErrConnectionNotFound       = errors.New("connection not found")
	ErrConnectionNotReady       = errors.New("connection not ready")
	ErrCreatePayloadFailed      = errors.New("create payload failed")
	ErrHandshakeFailed          = errors.New("handshake failed")
	ErrHandshakeTimeout         = errors.New("handshake timeout")
	ErrHandshakeCancelled       = errors.New("handshake cancelled")
	ErrPrepareChatSessionFailed = errors.New("prepare chat session failed")
)
