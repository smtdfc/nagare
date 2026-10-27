package host

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/plugin/internal"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/pkgs/helpers"
	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
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

			err := conn.Send(plugin_dtos.PluginToolCallEvent, plugin_dtos.PluginToolCallEventPayload{
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
				plugin_dtos.ReceivedChatMessageEvent,
				&plugin_dtos.ReceivedChatMessageEventPayload{
					SessionID: chunkPayload.SessionID,
					Message:   chunkJSON,
					ChannelID: chunkPayload.ChannelID,
				},
				chunkPayload.RequestID,
			)
		}
	}
}

func (h *PluginHost) forwardVectorIndexActions() {
	chAddVectorAction, unsubscribeAddVector := h.eventBus.Subscribe(event_bus.AddVectorEvent)
	defer unsubscribeAddVector()

	handleChan := helpers.MergeChannel(
		chAddVectorAction,
	)

	for {
		select {
		case <-h.stopChan:
			return
		case evt, ok := <-handleChan:
			if !ok {
				return
			}
			switch evt.GetEventType() {
			case event_bus.AddVectorEvent:
				payload, ok := evt.(*event_bus.AddVectorEventPayload)
				if !ok {
					continue
				}

				h.mu.RLock()
				conn, exist := h.specialGroups[VectorIndexPluginGroup]
				h.mu.RUnlock()

				if !exist {
					h.eventBus.Publish(context.Background(), event_bus.AddVectorErrorEvent, &event_bus.AddVectorErrorEventPayload{
						RequestID: payload.RequestID,
						Error:     custom_errors.ErrMissingVectorIndexPlugin,
					})

					continue
				}
				err := conn.Send(internal.AddVectorEvent, &internal.AddVectorEventPayload{
					RequestID: payload.RequestID,
					Vector:    payload.Vector,
					IDs:       payload.IDs,
				}, payload.RequestID)

				if err != nil {
					h.eventBus.Publish(context.Background(), event_bus.AddVectorErrorEvent, &event_bus.AddVectorErrorEventPayload{
						RequestID: payload.RequestID,
						Error:     err,
					})
				}
			}
		}
	}
}
