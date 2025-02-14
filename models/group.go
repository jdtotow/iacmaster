package models

import (
	"gorm.io/gorm"
)

// User Group structure
type UserGroup struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
	Uuid string `gorm:"primaryKey" json:"uuid"`
}

func (group UserGroup) GetName() string {
	return group.Name
}
func (group *UserGroup) SetName(name string) {
	group.Name = name
}

func (group *UserGroup) SetUuid(uuid string) {
	group.Uuid = uuid
}
