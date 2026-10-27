package event_bus

var (
	PluginToolCallEvent       EventType = "event:plugin:tool:call"
	PluginToolCallResultEvent EventType = "event:plugin:tool:result"
)

type PluginToolCallEventPayload struct {
	RequestID string
	PluginID  string
	Name      string
	Args      string
}

func (c *PluginToolCallEventPayload) GetEventType() EventType {
	return PluginToolCallEvent
}

type PluginToolCallResultEventPayload struct {
	RequestID string
	Result    string
	Error     string
}

func (c *PluginToolCallResultEventPayload) GetEventType() EventType {
	return PluginToolCallResultEvent
}
