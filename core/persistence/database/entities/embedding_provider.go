package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmbeddingProvider struct {
	gorm.Model
	ID         uuid.UUID `gorm:"primaryKey;" json:"id"`
	Name       string
	Compatible string
	Models     string
	BaseURL    string

	CredentialID uuid.UUID   `gorm:"type:uuid;not null" json:"credential_id"`
	Credential   *Credential `gorm:"foreignKey:CredentialID;references:ID" json:"credential"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *EmbeddingProvider) BeforeCreate(*gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
