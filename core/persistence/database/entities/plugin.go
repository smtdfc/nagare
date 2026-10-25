package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plugin struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey;" json:"id"`
	PackageName string    `gorm:"column:package_name;type:varchar(100);uniqueIndex;not null" json:"package_name"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Author      string    `gorm:"type:varchar(100)" json:"author"`
	Features    string    `gorm:"type:text" json:"features"`
	Version     string    `gorm:"type:varchar(50);not null" json:"version"`
	Bin         string    `gorm:"type:varchar(255)" json:"bin"`
	BinChecksum string    `gorm:"column:bin_checksum;type:varchar(64)" json:"bin_checksum"`
	// SignatureStatus is the archive verification result at install time:
	// "verified", "unverified", or "unsigned".
	SignatureStatus string `gorm:"column:signature_status;type:varchar(20);default:'unsigned'" json:"signature_status"`
	IsActive        bool   `gorm:"default:true" json:"is_active"`
	// LastStatus is the outcome of the most recent start or stop attempt:
	// "running", "stopped", "failed", or "unknown".
	LastStatus string `gorm:"column:last_status;type:varchar(20);default:'unknown'" json:"last_status"`

	// Sessions []Session `gorm:"foreignKey:PluginID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sessions,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Plugin) BeforeCreate(*gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
