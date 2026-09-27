package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DynamicTool struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex:idx_plugin_tool;not null"`
	Description string
	ArgsSchema  string
	PluginID    uuid.UUID `gorm:"uniqueIndex:idx_plugin_tool;not null" json:"plugin_id"`
	Plugin      *Plugin   `gorm:"foreignKey:PluginID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"plugin,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
