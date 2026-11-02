package cronjob

import (
	"github.com/smtdfc/nagare/core/cronjob"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func parseCronJobStatus(s string) cronjob.Status {
	return cronjob.MapStringToCronJobStatus(s)
}

func formatCronJobStatus(s cronjob.Status) string {
	return s.ToString()
}

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
// goverter:skipCopySameType
// goverter:useZeroValueOnPointerInconsistency
type ICronJobMapper interface {
	// goverter:map Status | parseCronJobStatus
	ToDomain(entity *entities.CronJob) *cronjob.CronJob

	// goverter:map Status | formatCronJobStatus
	ToEntity(domain *cronjob.CronJob) *entities.CronJob

	ToDomains(entities []*entities.CronJob) []*cronjob.CronJob
	ToEntities(domains []*cronjob.CronJob) []*entities.CronJob
}
