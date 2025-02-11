package models

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex"`
	Admin     int
	Variables []*EnvironmentVariable
	ProjectID uint
}

func CreateOrganization(name string) Organization {
	return Organization{
		Name: name,
	}
}
func (org *Organization) SetName(name string) {
	org.Name = name
}
func (org *Organization) SetAdmin(userID int) {
	org.Admin = userID
}
