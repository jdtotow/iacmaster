package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IaCArtifact struct {
	gorm.Model
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	ScmUrl     string    `json:"scm_url"`
	HomeFolder string    `json:"home_folder"`
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Project    Project
	ProjectID  uuid.UUID `json:"project_id"`
}

func (arti IaCArtifact) GetType() string {
	return arti.Type
}
func (arti IaCArtifact) GetName() string {
	return arti.Name
}
func (arti IaCArtifact) GetSCMUrl() string {
	return arti.ScmUrl
}
func (arti *IaCArtifact) SetType(_type string) {
	arti.Type = _type
}
func (arti *IaCArtifact) SetName(name string) {
	arti.Name = name
}
func (arti *IaCArtifact) SetSCMurl(url string) {
	arti.ScmUrl = url
}
