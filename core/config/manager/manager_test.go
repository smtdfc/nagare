package manager

import (
	"context"
	"log/slog"
	"testing"

	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestConfigManager_GetAndSet(t *testing.T) {
	// Setup in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	_ = db.AutoMigrate(&entities.KV{})

	l := &logger.BaseLogger{Logger: *slog.Default()}
	kvRepo := repositories.NewKVRepository(db, l)
	kvMapper := mappers.NewKVMapper()

	mgr := NewConfigManager(kvRepo, kvMapper, l)
	ctx := context.Background()

	// Initial get
	cfg, err := mgr.GetGeneralConfig(ctx)
	if err != nil {
		t.Fatalf("unexpected error on GetGeneralConfig: %v", err)
	}
	if cfg.DefaultLLMProvider != "" || cfg.DefaultLLMModel != "" {
		t.Errorf("expected empty initial config")
	}

	// Set config
	newCfg := &config.GeneralConfig{
		DefaultLLMProvider: "openai",
		DefaultLLMModel:    "gpt-4o",
	}
	err = mgr.SetGeneralConfig(ctx, newCfg)
	if err != nil {
		t.Fatalf("unexpected error on SetGeneralConfig: %v", err)
	}

	// Get updated config
	updated, err := mgr.GetGeneralConfig(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting updated config: %v", err)
	}
	if updated.DefaultLLMProvider != "openai" || updated.DefaultLLMModel != "gpt-4o" {
		t.Errorf("expected updated config, got %+v", updated)
	}
}
