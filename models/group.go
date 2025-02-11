package models

import (
	"gorm.io/gorm"
)

// User Group structure
type UserGroup struct {
	gorm.Model
	Name   string
	UserID uint
}

func CreateUserGroup(name string) UserGroup {
	return UserGroup{
		Name: name,
	}
}
func (group UserGroup) GetName() string {
	return group.Name
}
func (group *UserGroup) SetName(name string) {
	group.Name = name
}
