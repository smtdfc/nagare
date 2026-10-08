package mappers

import (
	"testing"

	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestKVMapper_ToDomains(t *testing.T) {
	// Verify converting slice of KV entities to domain map
	mapper := NewKVMapper()
	list := []*entities.KV{
		{Key: "theme", Value: "dark"},
		{Key: "lang", Value: "vi"},
	}

	result := mapper.ToDomains(list)
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	if result["theme"] != "dark" {
		t.Errorf("expected theme dark, got %s", result["theme"])
	}
	if result["lang"] != "vi" {
		t.Errorf("expected lang vi, got %s", result["lang"])
	}
}

func TestKVMapper_ToEntities(t *testing.T) {
	// Verify converting domain map to slice of KV entities
	mapper := NewKVMapper()
	domains := map[string]string{
		"debug": "true",
	}

	entitiesList := mapper.ToEntities(domains, "global")
	if len(entitiesList) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(entitiesList))
	}
	if entitiesList[0].Key != "debug" || entitiesList[0].Value != "true" || entitiesList[0].Scope != "global" {
		t.Errorf("expected key 'debug', value 'true', scope 'global', got %+v", entitiesList[0])
	}
}
