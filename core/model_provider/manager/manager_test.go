package manager

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/logger"
	model_provider_mapper "github.com/smtdfc/nagare/core/mappers/model_provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestModelProviderManager_GetAllAndGetByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	_ = db.AutoMigrate(&entities.ModelProvider{}, &entities.Credential{})

	l := &logger.BaseLogger{Logger: *slog.Default()}
	credRepo := repositories.NewCredentialRepository(db, l)
	modelMapper := model_provider_mapper.NewModelProviderMapper()
	modelRepo := repositories.NewModelProviderRepository(db, l)

	mgr := NewModelProviderManager(modelRepo, credRepo, modelMapper, l)
	if mgr == nil {
		t.Fatalf("expected non-nil manager")
	}

	ctx := context.Background()
	id := uuid.New()
	p := &entities.ModelProvider{
		ID:         id,
		Name:       "openai",
		Compatible: "OpenAI",
		ApiKey:     "sk-123",
		Models:     "gpt-4o",
	}
	if err := db.Create(p).Error; err != nil {
		t.Fatalf("failed to seed provider: %v", err)
	}

	all, err := mgr.GetAllProvider(ctx)
	if err != nil {
		t.Fatalf("GetAllProvider failed: %v", err)
	}
	if len(all) != 1 || all[0].ID != id {
		t.Fatalf("expected 1 provider with id %s, got %+v", id, all)
	}

	one, err := mgr.GetProviderByID(ctx, id.String())
	if err != nil {
		t.Fatalf("GetProviderByID failed: %v", err)
	}
	if one == nil || one.Name != "openai" {
		t.Fatalf("expected provider openai, got %+v", one)
	}

	if _, err := mgr.GetProviderByID(ctx, uuid.New().String()); err == nil {
		t.Fatalf("expected error for unknown id")
	}
}
