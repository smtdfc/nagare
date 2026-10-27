package event_bus

type SenderType string

const (
	User   SenderType = "user"
	Plugin SenderType = "plugin"
	System SenderType = "system"
)
