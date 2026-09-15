package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LLMProvider struct {
	gorm.Model
	ID         uuid.UUID `gorm:"primaryKey;" json:"id"`
	Name       string
	Compatible string
	ApiKey     string
	Models     string
	BaseURL    string

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *LLMProvider) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
