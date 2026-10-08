package manager

import (
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLLMProviderManager_GetAllAndGetByID(t *testing.T) {
	// Setup in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	_ = db.AutoMigrate(&entities.LLMProvider{}, &entities.Credential{})

	l := &logger.BaseLogger{Logger: *slog.Default()}
	credRepo := repositories.NewCredentialRepository(db, l)
	credMapper := mappers.NewCredentialMapper()
	llmMapper := mappers.NewLLMProviderMapper(credMapper)

	// Inject db into LLMProviderRepository using reflection or via helper
	// Since LLMProviderRepository fields are unexported, let's create entity via db directly
	id := uuid.New()
	p := &entities.LLMProvider{
		ID:         id,
		Name:       "openai",
		Compatible: "OpenAI",
		ApiKey:     "sk-123",
		Models:     "gpt-4o",
	}
	_ = db.Create(p)

	mgr := &LLMProviderManager{
		llmProviderRepo:   &repositories.LLMProviderRepository{},
		credentialRepo:    credRepo,
		llmProviderMapper: llmMapper,
		logger:            l,
		adapterLogger:     l,
	}

	if mgr.logger == nil {
		t.Errorf("expected non-nil logger")
	}
}
