package host

import (
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/pkgs/helpers"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (h *PluginHost) forwardPluginToolCalls() {
	channel, unsubscribe := h.eventBus.Subscribe(event_bus.PluginToolCallEvent)
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

			err := conn.Send(contracts.PluginToolCallEvent, contracts.PluginToolCallEventPayload{
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
	ch, unsubscribe := h.eventBus.Subscribe(event_bus.ChatMessageChunkEvent)
	defer unsubscribe()

	for {
		select {
		case <-h.stopChan:
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if evt.GetEventType() != event_bus.ChatMessageChunkEvent {
				continue
			}
			chunkPayload, ok := evt.(*event_bus.ChatMessageChunkEventPayload)
			if !ok {
				continue
			}

			var conn *PluginConnection
			if chunkPayload.SessionOwnerType == string(session.PLUGIN) ||
				chunkPayload.SessionOwnerType == string(event_bus.Plugin) {
				h.mu.RLock()
				c, exists := h.connections[chunkPayload.SessionOwnerID]
				h.mu.RUnlock()

				if !exists {
					continue
				}

				conn = c
			} else if c, exist := h.chatListeners[chunkPayload.SessionID]; exist {
				conn = c
			} else {
				continue
			}

			if conn == nil {
				continue
			}

			chunkJSON, err := helpers.MarshalJson(chunkPayload.Chunk)
			if err != nil {
				continue
			}

			_ = conn.Send(
				contracts.ReceivedChatMessageEvent,
				&contracts.ReceivedChatMessageEventPayload{
					Message:   chunkJSON,
					ChannelID: chunkPayload.ChannelID,
				},
				chunkPayload.RequestID,
			)
		}
	}
}
