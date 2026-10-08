package mappers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/task"
)

func TestTaskMapper(t *testing.T) {
	// Verify TaskMapper ToDomain and ToEntity
	mapper := NewTaskMapper()
	id := uuid.New()
	sessID := uuid.New()
	now := time.Now()

	entity := &entities.Task{
		ID:                 id,
		Name:               "Backup Job",
		IsActive:           true,
		Prompt:             "Do backup",
		SessionID:          sessID,
		Status:             "pending",
		TriggerBy:          "scheduled",
		TriggerByEventName: "event_a",
		RepeatType:         "daily",
		StartTime:          &now,
		NextRunTime:        &now,
	}

	domain := mapper.ToDomain(entity)
	if domain == nil {
		t.Fatalf("expected non-nil domain")
	}
	if domain.ID != id {
		t.Errorf("expected ID %s, got %s", id, domain.ID)
	}
	if domain.Status != task.Pending {
		t.Errorf("expected status Pending, got %s", domain.Status)
	}
	if domain.TriggerRule == nil || domain.TriggerRule.Repeat != task.Daily {
		t.Errorf("expected trigger repeat Daily")
	}

	entityBack := mapper.ToEntity(domain)
	if entityBack == nil || entityBack.Name != "Backup Job" || entityBack.Status != "pending" {
		t.Errorf("expected mapped entity back, got %+v", entityBack)
	}

	// Verify slice mapping
	domains := mapper.ToDomains([]*entities.Task{entity})
	if len(domains) != 1 {
		t.Errorf("expected 1 domain, got %d", len(domains))
	}

	entitiesList := mapper.ToEntities(domains)
	if len(entitiesList) != 1 {
		t.Errorf("expected 1 entity, got %d", len(entitiesList))
	}
}
