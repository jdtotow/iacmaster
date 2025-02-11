package models

import "gorm.io/gorm"

type IaCArtifact struct {
	gorm.Model
	Type   string
	Name   string
	ScmUrl string
}

func CreateIaCArtifact(_type, name, scm string) IaCArtifact {
	return IaCArtifact{
		Type:   _type,
		Name:   name,
		ScmUrl: scm,
	}
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
