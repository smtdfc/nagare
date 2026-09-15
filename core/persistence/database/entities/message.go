package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey;" json:"id"`
	MessageKind string    `gorm:"type:varchar(50);not null" json:"message_kind"`
	Content     string    `gorm:"type:text;not null" json:"content"`

	SessionID uuid.UUID `gorm:" index;not null" json:"session_id"`
	Session   *Session  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"session,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
