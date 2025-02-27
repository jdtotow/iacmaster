package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Group structure
type Group struct {
	gorm.Model
	Name  string    `gorm:"uniqueIndex"`
	ID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Users []User    `gorm:"many2many:user_groups;"`
}

func (group Group) GetName() string {
	return group.Name
}
func (group *Group) SetName(name string) {
	group.Name = name
}
