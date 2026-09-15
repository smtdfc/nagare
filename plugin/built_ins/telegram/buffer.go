package main

import (
	"strings"
	"sync"
)

// MessageBuffer stores accumulated text chunks for a session thread-safely.
type MessageBuffer struct {
	mu     sync.Mutex
	buffer strings.Builder
}

func (mb *MessageBuffer) Append(text string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.buffer.WriteString(text)
}

func (mb *MessageBuffer) String() string {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return mb.buffer.String()
}

func (tp *TelegramPlugin) getOrCreateBuffer(sessionID string) *MessageBuffer {
	val, _ := tp.messageBuffers.LoadOrStore(sessionID, &MessageBuffer{})
	return val.(*MessageBuffer)
}

func (tp *TelegramPlugin) clearBuffer(sessionID string) {
	tp.messageBuffers.Delete(sessionID)
}
