package models

import "gorm.io/gorm"

// Project structure
type Project struct {
	gorm.Model
	Name         string
	Parent       string
	Organization Organization
	Variables    []EnvironmentVariable
	Uuid         string
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
