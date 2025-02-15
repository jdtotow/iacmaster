package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CloudCredential struct {
	gorm.Model
	ID        uint
	Name      string                 `gorm:"uniqueIndex"`
	Type      DestinationCloud       `json:"destination_cloud"`
	Variables []EnvironmentVariable  `json:"variables"`
	Uuid      uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Settings  []IaCExecutionSettings `json:"settings"`
}

func (c *CloudCredential) SetType(_type string) {
	c.Type = DestinationCloud(_type)
}
func (c CloudCredential) HasVariable(name string) bool {
	for _, _var := range c.Variables {
		if _var.GetName() == name {
			return true
		}
	}
	return false
}
func (c *CloudCredential) AddVariable(_var EnvironmentVariable) {
	if !c.HasVariable(_var.Name) {
		c.Variables = append(c.Variables, _var)
	}
}
func (c CloudCredential) GetCloud(name DestinationCloud) []EnvironmentVariable {
	result := []EnvironmentVariable{}
	for _, _var := range c.Variables {
		if _var.Type == EnvVariableType(name+"_credential") {
			result = append(result, _var)
		}
	}
	return result
}
