package models

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
	Uuid string
}

func CreateOrganization(name string) Organization {
	return Organization{
		Name: name,
	}
}
func (org *Organization) SetName(name string) {
	org.Name = name
}
func (org *Organization) GetName() string {
	return org.Name
}
func (org *Organization) SetUuid(uuid string) {
	org.Uuid = uuid
}
