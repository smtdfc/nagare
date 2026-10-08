package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/pkgs/messages"
)

func TestMessageMapper(t *testing.T) {
	// Verify MessageMapper ToEntity and ToDomain
	mapper := NewMessageMapper()
	sessID := uuid.New()
	invokeID := uuid.New()

	txtMsg := messages.NewTextMessage(messages.USER, "Hello world")
	txtMsg.SetInvokeID(invokeID.String())

	entity, err := mapper.ToEntity(txtMsg, sessID.String())
	if err != nil {
		t.Fatalf("unexpected error on ToEntity: %v", err)
	}
	if entity == nil {
		t.Fatalf("expected non-nil entity")
	}
	if entity.SessionID != sessID {
		t.Errorf("expected session ID %s, got %s", sessID, entity.SessionID)
	}
	if entity.InvokeID != invokeID {
		t.Errorf("expected invoke ID %s, got %s", invokeID, entity.InvokeID)
	}

	// Map back to domain
	domainBack, err := mapper.ToDomain(entity)
	if err != nil {
		t.Fatalf("unexpected error on ToDomain: %v", err)
	}
	if domainBack == nil {
		t.Fatalf("expected non-nil domain back")
	}
	if domainBack.GetMessageType() != messages.TextMessageType {
		t.Errorf("expected TextMessageType, got %s", domainBack.GetMessageType())
	}

	// Verify slice mapping
	entitiesList, err := mapper.ToEntities([]messages.Message{txtMsg}, sessID.String())
	if err != nil || len(entitiesList) != 1 {
		t.Errorf("expected 1 entity from slice, got %d (err: %v)", len(entitiesList), err)
	}

	domainsList, err := mapper.ToDomains(entitiesList)
	if err != nil || len(domainsList) != 1 {
		t.Errorf("expected 1 domain from slice, got %d (err: %v)", len(domainsList), err)
	}

	// Verify nil handling
	eNil, err := mapper.ToEntity(nil, sessID.String())
	if err != nil || eNil != nil {
		t.Errorf("expected nil entity for nil message")
	}

	dNil, err := mapper.ToDomain(nil)
	if err != nil || dNil != nil {
		t.Errorf("expected nil domain for nil entity")
	}

	// Verify unknown message kind error
	unknownEntity := &entities.Message{
		MessageKind: "unknown_kind",
		Content:     "{}",
	}
	_, err = mapper.ToDomain(unknownEntity)
	if err == nil {
		t.Errorf("expected error for unknown message kind")
	}
}
