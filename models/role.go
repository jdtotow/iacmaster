package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name   string `json:"name"`
	UserID uuid.UUID
	User   User
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

func (r Role) GetName() string {
	return r.Name
}
func (r *Role) SetName(name string) {
	r.Name = name
}
