package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID       uint
	Name     string `json:"name"`
	UserUuid uuid.UUID
	User     User      `gorm:"foreignKey:UserUuid"`
	Uuid     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

func (r Role) GetName() string {
	return r.Name
}
func (r *Role) SetName(name string) {
	r.Name = name
}
