package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:char(36);primaryKey;"`

	MessageType string
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	SessionID uuid.UUID `gorm:"type:char(36);"`
	Session   Session
}

func (s *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
