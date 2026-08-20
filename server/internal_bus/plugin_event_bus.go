package internal_bus

import "github.com/smtdfc/nagare/shared/bus"

type PluginEventBus struct {
	bus *bus.EventBus
}

func (p *PluginEventBus) Subscribe(topic string) (<-chan bus.Event, func()) {
	return p.bus.Subscribe(topic, 1)
}

// @Injectable
func NewPluginEventBus() *PluginEventBus {
	return &PluginEventBus{
		bus: bus.NewEventBus(),
	}
}
