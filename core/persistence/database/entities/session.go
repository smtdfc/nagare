package entities

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primaryKey;" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	ChannelID string    `gorm:"type:varchar(255);index:idx_channel,composite:owner" json:"channel_id"`
	OwnerID   uuid.UUID `gorm:"type:uuid;index:idx_owner,composite:owner" json:"owner_id"`
	OwnerType string    `gorm:"type:varchar(50);index:idx_owner,composite:owner" json:"owner_type"`
	IsArchive bool      `gorm:"default:false" json:"is_archive"`

	Messages []*Message `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"messages,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
