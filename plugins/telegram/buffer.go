package main

import (
	"strings"
	"sync"
)

// MessageBuffer stores accumulated text chunks for a channel thread-safely.
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

func (tp *TelegramPlugin) getOrCreateBuffer(channelID string) *MessageBuffer {
	val, _ := tp.messageBuffers.LoadOrStore(channelID, &MessageBuffer{})
	return val.(*MessageBuffer)
}

func (tp *TelegramPlugin) clearBuffer(channelID string) {
	tp.messageBuffers.Delete(channelID)
}
