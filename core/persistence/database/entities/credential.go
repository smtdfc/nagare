package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	ID uuid.UUID `gorm:"primaryKey;" json:"id"`

	Name   string `gorm:"type:varchar(255);not null"`
	ApiKey string `gorm:"type:varchar(255);not null"`

	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

func (p *Credential) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
