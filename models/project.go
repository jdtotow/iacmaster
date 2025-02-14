package models

import "gorm.io/gorm"

// Project structure
type Project struct {
	gorm.Model
	Name         string                `json:"name"`
	Parent       string                `json:"parent"`
	Organization Organization          `json:"organization"`
	VariableUuid string                `json:"variable_uuid"`
	Variables    []EnvironmentVariable `json:"variables" gorm:"foreignKey:Uuid;references:VariableUuid"`
	Uuid         string                `gorm:"primaryKey" json:"uuid"`
}

func (project Project) GetName() string {
	return project.Name
}
func (project Project) GetParent() string {
	return project.Parent
}
func (project Project) GetOrganization() Organization {
	return project.Organization
}
func (project *Project) SetUuid(uuid string) {
	project.Uuid = uuid
}
