package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project structure
type Project struct {
	gorm.Model
	Name            string            `json:"name"`
	Parent          string            `json:"parent"`
	OrganizationID  uuid.UUID         `json:"organization_id"`
	Organization    Organization      `json:"organization"`
	Variables       map[string]string `json:"variables" gorm:"type:jsonb"`
	ID              uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Environment     []Environment     `json:"environments"`
	CloudCredential []CloudCredential `json:"credentials"`
	Token           []Token           `json:"tokens"`
	IaCArtifact     []IaCArtifact
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
