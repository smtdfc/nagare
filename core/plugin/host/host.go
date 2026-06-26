package host

import (
	"encoding/json"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/plugin/features"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/plugin-sdk/host"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

func StartHost(host *host.Host, pluginMgr *manager.PluginManager, chatMgr *features.ChatChannelManager, pool *agent.AgentPool, sessionMgr *agent.SessionManager) {
	host.Handler(shared.REGISTER_CHAT_CHANNEL, func(msg shared.Message) {
		var payload shared.RegisterChatChannelPayload
		json.Unmarshal(msg.Payload, &payload)
		a := pool.GetOrNew()
		a.State.ExtendHistory(
			sessionMgr.GetHistory(payload.ID, agent.NAGARE_LIST_MESSAGE_SIZE_LIMIT),
		)
		chatMgr.Register(&features.ChatChannel{
			Id:         payload.ID,
			Agent:      a,
			SessionMgr: sessionMgr,
			SessionID:  payload.ID,
			CleanUp: func() {
				sessionMgr.SaveHistory(payload.ID, a.State.History)
				pool.Put(a)
			},
		})

		host.Send(msg.PluginID, shared.REGISTER_CHAT_CHANNEL_SUCCESS, shared.RegisterChatChannelSuccessPayload{
			ID: payload.ID,
		})
	})

	host.Handler(shared.HANDLE_CHAT_MESSAGE, func(msg shared.Message) {
		var payload shared.HandleChatMessagePayload
		json.Unmarshal(msg.Payload, &payload)
		chatMgr.Handle(&payload, msg.PluginID)
	})

	host.Start()
}
