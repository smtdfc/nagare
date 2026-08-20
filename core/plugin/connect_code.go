package plugin

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/domains"
)

type pluginCodeEntry struct {
	plugin    *domains.PluginInfo
	expiresAt time.Time
}

type ConnectCodeManager struct {
	mu    sync.RWMutex
	codes map[string]pluginCodeEntry
}

func (m *ConnectCodeManager) GenerateConnectCode(plugin *domains.PluginInfo, ttl time.Duration) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	code := uuid.New().String()
	m.codes[code] = pluginCodeEntry{
		plugin:    plugin,
		expiresAt: time.Now().Add(ttl),
	}

	return code
}

func (m *ConnectCodeManager) HasConnectCode(code string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.codes[code]
	if !ok {
		return false
	}

	if time.Now().After(entry.expiresAt) {
		return false
	}

	return ok
}

func (m *ConnectCodeManager) GetPluginFromCode(code string) *domains.PluginInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.codes[code]
	if !ok {
		return nil
	}

	defer delete(m.codes, code)

	if time.Now().After(entry.expiresAt) {
		return nil
	}

	return entry.plugin
}

// @Injectable
func NewConnectCodeManager() *ConnectCodeManager {
	manager := &ConnectCodeManager{
		codes: make(map[string]pluginCodeEntry),
	}

	return manager
}
