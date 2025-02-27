package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CloudCredential struct {
	gorm.Model
	Name                 string                 `json:"name" gorm:"uniqueIndex"`
	Type                 DestinationCloud       `json:"destination_cloud"`
	Variables            map[string]string      `json:"variables" gorm:"type:jsonb;serializer:json"`
	ID                   uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	IaCExecutionSettings []IaCExecutionSettings `json:"settings"`
	Project              Project                `json:"project"`
	ProjectID            uuid.UUID              `json:"project_id"`
}

func (c *CloudCredential) SetType(_type string) {
	c.Type = DestinationCloud(_type)
}
func (c CloudCredential) HasVariable(name string) bool {
	for _name := range c.Variables {
		if _name == name {
			return true
		}
	}
	return false
}
func (c *CloudCredential) AddVariable(name, _value string) {
	c.Variables[name] = _value
}
