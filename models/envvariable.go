package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnvVariableType string

const AzureCredential EnvVariableType = "azure_credential"
const AWSCredential EnvVariableType = "aws_credential"
const GCPCredential EnvVariableType = "gcp_credential"
const GENERAL EnvVariableType = "general"

type EnvironmentVariable struct {
	gorm.Model
	ID       uint
	Type     EnvVariableType `json:"type"`
	Name     string          `json:"name"`
	Value    string          `json:"value"`
	Uuid     uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Projects []Project       `json:"projects" gorm:"many2many:project_variables;"`
}

func CreateEnvironmentVariable(name, value string) EnvironmentVariable {
	return EnvironmentVariable{
		Name:  name,
		Value: value,
	}
}
func (env EnvironmentVariable) GetName() string {
	return env.Name
}
func (env EnvironmentVariable) GetValue() string {
	return env.Value
}
func (env *EnvironmentVariable) SetName(name string) {
	env.Name = name
}
func (env *EnvironmentVariable) SetValue(value string) {
	env.Value = value
}
