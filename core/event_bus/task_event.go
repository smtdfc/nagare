package event_bus

var (
	RefreshTaskEvent EventType = "event:task:refresh"
)

type TaskInfo struct {
	ID                    string `json:"id"`
	SendResultIntoSession string `json:"sendResultIntoSession"`
}

type RefreshTaskEventPayload struct {
}

func (c *RefreshTaskEventPayload) GetEventType() EventType {
	return RefreshTaskEvent
}
