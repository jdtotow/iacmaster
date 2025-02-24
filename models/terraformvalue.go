package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TerraformValue struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Value []byte    `json:"value" gorm:"type:jsonb;serializer:json"`
}
