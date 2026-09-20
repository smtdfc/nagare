package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey;" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
	Prompt   string    `gorm:"type:text;" json:"prompt"`
	// Sessions []Session `gorm:"foreignKey:PluginID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sessions,omitempty"`

	Status             string `gorm:"type:text;" json:"status"`
	TriggerBy          string `gorm:"type:text;" json:"trigger_by"`
	TriggerByEventName string `gorm:"type:text;" json:"trigger_by_event_name"`

	RepeatType  string     `gorm:"type:text; default:'no_repeat'" json:"trigger_repeat_type" `
	StartTime   *time.Time `gorm:"type:timestamp;" json:"start_time"`
	EndTime     *time.Time `gorm:"type:timestamp;" json:"end_time"`
	NextRunTime *time.Time `gorm:"type:timestamp;index" json:"next_run_time"`

	SessionID uuid.UUID `gorm:" index;not null" json:"session_id"`
	Session   *Session  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"session,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
