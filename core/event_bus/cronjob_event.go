package event_bus

var (
	RefreshCronJobEvent EventType = "event:cronjob:refresh"
)

type CronJobInfo struct {
	ID                    string `json:"id"`
	SendResultIntoSession string `json:"sendResultIntoSession"`
}

type RefreshCronJobEventPayload struct{}

func (c *RefreshCronJobEventPayload) GetEventType() EventType {
	return RefreshCronJobEvent
}
