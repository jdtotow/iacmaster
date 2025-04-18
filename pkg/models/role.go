package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name  string    `json:"name"`
	Users []User    `json:"users" gorm:"many2many:user_roles;"`
	ID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

func (r Role) GetName() string {
	return r.Name
}
func (r *Role) SetName(name string) {
	r.Name = name
}
