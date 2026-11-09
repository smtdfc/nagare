package external_module

import (
	"sync"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/logger"
)

type ExternalModuleManager struct {
	mu              sync.RWMutex
	logger          *logger.BaseLogger
	connectCodes    map[string]string
	ExternalModules map[string]ExternalModule
}

func (e *ExternalModuleManager) CreateConnectCode(moduleID string) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	code := uuid.New().String()
	e.connectCodes[code] = moduleID
	return code
}

func (e *ExternalModuleManager) IsValidConnect(moduleID string, connectCode string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if _, ok := e.ExternalModules[moduleID]; !ok {
		return false
	}

	if m, ok := e.connectCodes[connectCode]; !ok || moduleID != m {
		return false
	}

	return true
}

func (e *ExternalModuleManager) RevokeConnectCode(connectCode string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.connectCodes, connectCode)
}

// @Injectable
func NewExternalModuleManager(
	logger *logger.BaseLogger,
) *ExternalModuleManager {
	return &ExternalModuleManager{
		connectCodes: make(map[string]string),
		ExternalModules: map[string]ExternalModule{
			"nagare_vector":       NewExternalModule("nagare_vector"),
			"nagare_computer_use": NewExternalModule("nagare_computer_use"),
		},
		logger: logger.With("module", "external-module-manager"),
	}
}
