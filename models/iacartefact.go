package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IaCArtifact struct {
	gorm.Model
	ID     uint
	Type   string    `json:"type"`
	Name   string    `json:"name"`
	ScmUrl string    `json:"scm_url"`
	Uuid   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
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
