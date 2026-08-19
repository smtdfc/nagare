package handlers

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/ws"
)

type PluginConnectionEntry struct {
	Conn *ws.WsInstance
	Info *domains.PluginInfo
}

type PluginConnectionRegistry struct {
	mu      sync.RWMutex
	clients map[string]*PluginConnectionEntry
}

func (r *PluginConnectionRegistry) Register(pluginID string, conn *ws.WsInstance, info *domains.PluginInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entry, exists := r.clients[pluginID]; exists {
		entry.Conn.Close()
	}

	r.clients[pluginID] = &PluginConnectionEntry{
		Conn: conn,
		Info: info,
	}
}

func (r *PluginConnectionRegistry) Unregister(pluginID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, pluginID)
}

func (r *PluginConnectionRegistry) GetConn(pluginID string) (*ws.WsInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, ok := r.clients[pluginID]
	if !ok {
		return nil, fmt.Errorf("plugin %s is not connected", pluginID)
	}
	return entry.Conn, nil
}

func (r *PluginConnectionRegistry) GetPluginInfo(pluginID string) (*domains.PluginInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, ok := r.clients[pluginID]
	if !ok {
		return nil, fmt.Errorf("plugin %s is not connected", pluginID)
	}
	return entry.Info, nil
}

func (r *PluginConnectionRegistry) Broadcast(event dto.WsEvent, payload any) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, entry := range r.clients {
		_ = ws.SendMessage(entry.Conn, event, payload)
	}
}

// @Injectable
func NewPluginConnectionRegistry() *PluginConnectionRegistry {
	return &PluginConnectionRegistry{
		clients: make(map[string]*PluginConnectionEntry),
	}
}

type PluginHandler struct {
	pluginConnRegistry *PluginConnectionRegistry
}

func (h *PluginHandler) Register(i *ws.WsInstance, message *dto.WsMessage[any]) {
	payload, err := ws.GetPayload[dto.PluginRegisterEvent](message)
	if err != nil {
		ws.SendMessage(i, dto.WS_PLUGIN_REGISTER_FAILED, dto.PluginRegisterFailedEvent{
			Cause: "payload error: ",
		})
	}

	if !global.GlobalPluginMgr.HasConnectCode(payload.Code) {
		ws.SendMessage(i, dto.WS_PLUGIN_REGISTER_FAILED, dto.PluginRegisterFailedEvent{
			Cause: "invalid connect code",
		})
		return
	}
	pluginInfo := global.GlobalPluginMgr.GetPluginByConnectCode(payload.Code)
	i.Auth = &dto.AuthPayload{
		ID:     pluginInfo.ID,
		Target: "plugin",
	}

	h.pluginConnRegistry.Register(pluginInfo.ID, i, pluginInfo)
	ws.SendMessage(i, dto.WS_PLUGIN_REGISTER_SUCCESS, dto.PluginRegisterSuccessEvent{})
}

// @Injectable
func NewPluginHandler(pluginConnRegistry *PluginConnectionRegistry) *PluginHandler {
	return &PluginHandler{
		pluginConnRegistry: pluginConnRegistry,
	}
}
