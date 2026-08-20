package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plugin struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:char(36);primaryKey;"`

	PluginID   string
	Name       string
	Active     bool   `gorm:"default:true"`
	ApiVersion string `gorm:"default:v1"`
	Author     string `gorm:"default:unknown"`
	Version    string `gorm:"default:v1"`
	Bin        string `gorm:"default:unknown"`
	Features   string `gorm:"default:basic"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s *Plugin) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
