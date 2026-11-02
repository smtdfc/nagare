package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CronJob struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey;" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
	Prompt   string    `gorm:"type:text;" json:"prompt"`

	Status   string `gorm:"type:text;" json:"status"`
	CronExpr string `gorm:"type:text;not null" json:"cron_expr"`

	NextRunTime *time.Time `gorm:"type:timestamp;index" json:"next_run_time"`
	LastRunTime *time.Time `gorm:"type:timestamp;" json:"last_run_time"`

	SessionID uuid.UUID `gorm:"index;not null" json:"session_id"`
	Session   *Session  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"session,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (c *CronJob) BeforeCreate(*gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
