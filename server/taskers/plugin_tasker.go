package taskers

import (
	"github.com/smtdfc/nagare/server/internal_bus"
	"github.com/smtdfc/nagare/server/ws/handlers"
)

type PluginTasker struct {
}

// @Injectable
func InitPluginTasker(registry *handlers.PluginConnectionRegistry, bus *internal_bus.PluginEventBus) *PluginTasker {
	// go (func() {
	// 	channel, unsub := bus.Subscribe("")
	// })()

	return &PluginTasker{}
}
