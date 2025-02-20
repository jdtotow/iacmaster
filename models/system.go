package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IaCSystem struct {
	gorm.Model
	Name string
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}
