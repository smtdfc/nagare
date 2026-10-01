package host

import (
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

func (h *PluginHost) forwardPluginToolCalls() {
	channel, unsubscribe := h.chatEventBus.Subscribe(event_bus.PluginToolCallEvent)
	defer unsubscribe()

	for {
		select {
		case <-h.stopChan:
			return
		case payload, ok := <-channel:
			if !ok {
				return
			}
			call, ok := payload.(*event_bus.PluginToolCallEventPayload)
			if !ok {
				continue
			}

			h.mu.RLock()
			conn, exists := h.connections[call.PluginID]
			h.mu.RUnlock()

			if !exists {
				h.logger.Error("Plugin connection not found for tool call", "pluginID", call.PluginID, "name", call.Name, "requestID", call.RequestID)
				continue
			}

			err := conn.Send(plugin.PluginToolCallEvent, plugin.PluginToolCallEventPayload{
				Name: call.Name,
				Args: call.Args,
			}, call.RequestID)
			if err != nil {
				h.logger.Error("Failed to forward plugin tool call", "error", err, "pluginID", call.PluginID, "requestID", call.RequestID)
				continue
			}
			h.logger.Info("Forwarded plugin tool call", "pluginID", call.PluginID, "requestID", call.RequestID)
		}
	}
}

func (h *PluginHost) forwardChatChunks() {
	ch, unsubscribe := h.chatEventBus.Subscribe(event_bus.ChunkEvent)
	defer unsubscribe()

	for {
		select {
		case <-h.stopChan:
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if evt.GetEventType() != event_bus.ChunkEvent {
				continue
			}
			chunkPayload, ok := evt.(*event_bus.ChatChunkEventPayload)
			if !ok {
				continue
			}

			if chunkPayload.SessionOwnerType != string(session.PLUGIN) && chunkPayload.SessionOwnerType != string(event_bus.Plugin) {
				continue
			}

			h.mu.RLock()
			conn, exists := h.connections[chunkPayload.SessionOwnerID]
			h.mu.RUnlock()

			if !exists {
				continue
			}

			chunkJSON, err := helpers.MarshalJson(chunkPayload.Chunk)
			if err != nil {
				continue
			}

			_ = conn.Send(
				websocket.ReceivedChatMessageEvent,
				&websocket.ReceivedChatMessageEventPayload{
					SessionID: chunkPayload.SessionID,
					Message:   chunkJSON,
					ChannelID: chunkPayload.ChannelID,
				},
				chunkPayload.RequestID,
			)
		}
	}
}
