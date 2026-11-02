package cronjob

import (
	"testing"
	"time"

	"github.com/google/uuid"
	cronjob_domain "github.com/smtdfc/nagare/core/cronjob"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestCronJobMapper(t *testing.T) {
	mapper := NewCronJobMapper()

	now := time.Now()
	id := uuid.New()
	sessionID := uuid.New()

	entity := &entities.CronJob{
		ID:          id,
		Name:        "Daily sync",
		SessionID:   sessionID,
		IsActive:    true,
		Prompt:      "Sync data",
		Status:      "pending",
		CronExpr:    "0 12 * * *",
		NextRunTime: &now,
	}

	domain := mapper.ToDomain(entity)
	if domain == nil {
		t.Fatal("expected non-nil domain")
	}
	if domain.ID != id || domain.Name != "Daily sync" || domain.Status != cronjob_domain.Pending {
		t.Errorf("domain mapping mismatch: %+v", domain)
	}

	convertedEntity := mapper.ToEntity(domain)
	if convertedEntity == nil {
		t.Fatal("expected non-nil entity")
	}
	if convertedEntity.ID != id || convertedEntity.Name != "Daily sync" || convertedEntity.Status != "pending" {
		t.Errorf("entity mapping mismatch: %+v", convertedEntity)
	}
}
