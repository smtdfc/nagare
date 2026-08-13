package models

import (
	"errors"

	"gorm.io/gorm"
)

type KV struct {
	gorm.Model
	Target string `gorm:"type:text;"`
	Key    string `gorm:"type:text;"`
	Value  string `gorm:"type:text;"`
}

func (s *KV) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Key == "" {
		return errors.New("key is required")
	}
	if s.Value == "" {
		return errors.New("value is required")
	}
	return nil
}
