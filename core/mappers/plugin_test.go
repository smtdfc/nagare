package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/plugin"
)

func TestPluginMapper_ToDomain(t *testing.T) {
	// Verify PluginMapper ToDomain
	mapper := NewPluginMapper()
	id := uuid.New()

	entity := &entities.Plugin{
		ID:          id,
		PackageName: "com.example.test",
		Name:        "TestPlugin",
		Author:      "Author",
		Features:    "chat,plugin_tool",
		Version:     "1.0.0",
		Bin:         "test.bin",
		IsActive:    true,
	}

	domain := mapper.ToDomain(entity)
	if domain == nil {
		t.Fatalf("expected non-nil domain")
	}
	if domain.ID != id {
		t.Errorf("expected ID %s, got %s", id, domain.ID)
	}
	if len(domain.Features) != 2 {
		t.Errorf("expected 2 features, got %d", len(domain.Features))
	}

	// Verify nil handling
	if mapper.ToDomain(nil) != nil {
		t.Errorf("expected nil for nil entity")
	}
}

func TestPluginMapper_ToEntity(t *testing.T) {
	// Verify PluginMapper ToEntity
	mapper := NewPluginMapper()
	id := uuid.New()

	domain := &plugin.Plugin{
		ID:          id,
		PackageName: "com.example.test",
		Name:        "TestPlugin",
		Author:      "Author",
		Features:    []plugin.Feature{plugin.ChatFeature},
		Version:     "1.0.0",
		Bin:         "test.bin",
		IsActive:    true,
	}

	entity := mapper.ToEntity(domain)
	if entity == nil {
		t.Fatalf("expected non-nil entity")
	}
	if entity.Features != "chat" {
		t.Errorf("expected features 'chat', got '%s'", entity.Features)
	}
}

func TestPluginMapper_ToDomains(t *testing.T) {
	// Verify ToDomains batch mapping
	mapper := NewPluginMapper()
	list := []*entities.Plugin{
		{ID: uuid.New(), Name: "P1"},
		{ID: uuid.New(), Name: "P2"},
	}

	domains := mapper.ToDomains(list)
	if len(domains) != 2 {
		t.Errorf("expected 2 domains, got %d", len(domains))
	}
}
