package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name      string            `gorm:"uniqueIndex" json:"name"`
	ID        uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Variables map[string]string `json:"variables" gorm:"type:jsonb"`
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
