package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LLMProvider struct {
	gorm.Model

	ID              uuid.UUID `gorm:"type:char(36);primaryKey;"`
	Compatible      string    `gorm:"type:text;"`
	Name            string    `gorm:"type:text;"`
	BaseURL         string    `gorm:"type:text;"`
	APIKey          string    `gorm:"type:text;"`
	IsEnable        bool      `gorm:"type:boolean;"`
	ModelName       string    `gorm:"type:text;"`
	AvailableModels string    `gorm:"type:text;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *LLMProvider) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
