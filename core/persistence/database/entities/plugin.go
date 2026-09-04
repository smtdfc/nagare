package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plugin struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey;" json:"id"`
	PluginID string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"plugin_id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	Author   string    `gorm:"type:varchar(100)" json:"author"`
	Features string    `gorm:"type:text" json:"features"`
	Version  string    `gorm:"type:varchar(50);not null" json:"version"`
	Bin      string    `gorm:"type:varchar(255)" json:"bin"`
	IsActive bool      `gorm:"default:true" json:"is_active"`

	// Sessions []Session `gorm:"foreignKey:PluginID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sessions,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Plugin) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
