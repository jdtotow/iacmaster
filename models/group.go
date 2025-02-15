package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Group structure
type UserGroup struct {
	gorm.Model
	ID    uint
	Name  string    `gorm:"uniqueIndex"`
	Uuid  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Users []User    `gorm:"many2many:user_groups;"`
}

func (group UserGroup) GetName() string {
	return group.Name
}
func (group *UserGroup) SetName(name string) {
	group.Name = name
}
