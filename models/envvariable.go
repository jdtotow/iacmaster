package models

import "gorm.io/gorm"

type EnvVariableType string

const AzureCredential EnvVariableType = "azure_credential"
const AWSCredential EnvVariableType = "aws_credential"
const GCPCredential EnvVariableType = "gcp_credential"
const GENERAL EnvVariableType = "general"

type EnvironmentVariable struct {
	gorm.Model
	Type  EnvVariableType `json:"type"`
	Name  string          `json:"name"`
	Value string          `json:"value"`
	Uuid  string          `gorm:"primaryKey" json:"uuid"`
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
func (env *EnvironmentVariable) SetUuid(uuid string) {
	env.Uuid = uuid
}
