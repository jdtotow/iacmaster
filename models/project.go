package models

import "gorm.io/gorm"

// Project structure
type Project struct {
	gorm.Model
	Name           string
	Parent         string
	OrganizationID string
	Variables      []*EnvironmentVariable
}

func CreateProject(name, parent string, org string) Project {
	return Project{
		Name:           name,
		Parent:         parent,
		OrganizationID: org,
	}
}
func (project Project) GetName() string {
	return project.Name
}
func (project Project) GetParent() string {
	return project.Parent
}
func (project Project) GetOrganization() string {
	return project.OrganizationID
}
