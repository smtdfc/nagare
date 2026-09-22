package chat

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/messages"

	websocket_dtos "github.com/smtdfc/nagare/dtos/websocket"
)

type SessionJob struct {
	Chunk            messages.Message
	RequestID        string
	ChannelID        string
	SessionOwnerID   string
	SessionOwnerType string
}

type ChatWorker struct {
	mu           sync.RWMutex
	chatEventBus *event_bus.CoreEventBus
	logger       *logger.BaseLogger
	ws           *websocket.Coordinator

	sessionChans map[string]chan SessionJob
}

func (c *ChatWorker) HandleChunkMessage(sessionID string, job *SessionJob) {
	chunkJson, _ := helpers.MarshalJson(job.Chunk)
	err := websocket.BroadcastToRoom(
		c.ws,
		fmt.Sprintf("session:%s", sessionID),
		websocket_dtos.ReceivedChatMessageEvent,
		&websocket_dtos.ReceivedChatMessageEventPayload{
			SessionID: sessionID,
			Message:   chunkJson,
			ChannelID: job.ChannelID,
		},
		job.RequestID,
		nil,
	)
	if err != nil {
		c.logger.Error("Failed to broadcast chunk event: ", "requestID", job.RequestID, "error", err)
		return
	}

	if job.SessionOwnerID != "" && job.SessionOwnerType == string(session.PLUGIN) {
		err := websocket.BroadcastToRoom(
			c.ws,
			fmt.Sprintf("plugin:%s:chat", job.SessionOwnerID),
			websocket_dtos.ReceivedChatMessageEvent,
			&websocket_dtos.ReceivedChatMessageEventPayload{
				SessionID: sessionID,
				Message:   chunkJson,
				ChannelID: job.ChannelID,
			},
			job.RequestID,
			nil,
		)
		if err != nil {
			c.logger.Error("Failed to broadcast chunk event: ", "requestID", job.RequestID, "error", err)
			return
		}
	}
}

func (c *ChatWorker) runHandleSessionMessageWorker(sessionID string, ch chan SessionJob) {
	defer func() {
		c.mu.Lock()
		delete(c.sessionChans, sessionID)
		c.mu.Unlock()
		close(ch)
		c.logger.Debug("Session worker terminated and cleaned up", "session_id", sessionID)
	}()

	for job := range ch {
		c.HandleChunkMessage(sessionID, &job)
	}
}

func (c *ChatWorker) Do() {
	go func() {
		ch, unsubscribe := c.chatEventBus.Subscribe(event_bus.ChunkEvent)
		defer unsubscribe()

		for evt := range ch {
			switch evt.GetEventType() {
			case event_bus.ChunkEvent:
				chunkPayload, ok := evt.(*event_bus.ChatChunkEventPayload)
				if !ok {
					continue
				}

				sessionID := chunkPayload.SessionID

				c.mu.Lock()
				sessionChan, exists := c.sessionChans[sessionID]
				if !exists {
					sessionChan = make(chan SessionJob, 200)
					c.sessionChans[sessionID] = sessionChan

					go c.runHandleSessionMessageWorker(sessionID, sessionChan)
				}
				c.mu.Unlock()

				sessionChan <- SessionJob{
					Chunk:            chunkPayload.Chunk,
					RequestID:        chunkPayload.RequestID,
					ChannelID:        chunkPayload.ChannelID,
					SessionOwnerType: chunkPayload.SessionOwnerType,
					SessionOwnerID:   chunkPayload.SessionOwnerID,
				}
			}
		}
	}()
}

// @Injectable
func NewWorker(logger *logger.BaseLogger, chatEventBus *event_bus.CoreEventBus, ws *websocket.Coordinator) *ChatWorker {
	return &ChatWorker{
		chatEventBus: chatEventBus,
		mu:           sync.RWMutex{},
		ws:           ws,
		logger:       logger.With("worker", "gateway:chat:worker"),
		sessionChans: make(map[string]chan SessionJob),
	}
}
