package entities

import (
	"time"

	"gorm.io/gorm"
)

type KV struct {
	gorm.Model
	Value string
	Key   string `gorm:"uniqueIndex:idx_key_scope;not null"`
	Scope string `gorm:"uniqueIndex:idx_key_scope;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
