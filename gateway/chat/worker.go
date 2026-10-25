package chat

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/messages"

	websocket_dtos "github.com/smtdfc/nagare/dtos/websocket"
)

type SessionJob struct {
	Chunk            messages.Message
	InvokeID         string
	ChannelID        string
	SessionOwnerID   string
	SessionOwnerType string
}

type Worker struct {
	chatEventBus *event_bus.CoreEventBus
	logger       *logger.BaseLogger
	ws           *websocket.Coordinator
	doOnce       sync.Once
}

func (c *Worker) HandleChunkMessage(sessionID string, job *SessionJob) {
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
		job.InvokeID,
		nil,
	)
	if err != nil {
		c.logger.Warn("Failed to broadcast chunk event: ", "requestID", job.InvokeID, "error", err)
		return
	}
}

func (c *Worker) Do() {
	c.doOnce.Do(func() {
		go func() {
			ch, unsubscribe := c.chatEventBus.Subscribe(event_bus.ChunkEvent)
			defer unsubscribe()

			// Broadcast on the subscription goroutine. A per-session worker
			// queue used to be created here, but it was never cleaned up and
			// leaked a goroutine and channel for every session that streamed.
			for evt := range ch {
				chunkPayload, ok := evt.(*event_bus.ChatChunkEventPayload)
				if !ok || chunkPayload.SenderType == event_bus.Plugin {
					continue
				}

				c.HandleChunkMessage(chunkPayload.SessionID, &SessionJob{
					Chunk:            chunkPayload.Chunk,
					InvokeID:         chunkPayload.RequestID,
					ChannelID:        chunkPayload.ChannelID,
					SessionOwnerType: chunkPayload.SessionOwnerType,
					SessionOwnerID:   chunkPayload.SessionOwnerID,
				})
			}
		}()
	})
}

// @Injectable
func NewWorker(logger *logger.BaseLogger, chatEventBus *event_bus.CoreEventBus, ws *websocket.Coordinator) *Worker {
	return &Worker{
		chatEventBus: chatEventBus,
		ws:           ws,
		logger:       logger.With("worker", "gateway:chat:worker"),
	}
}
