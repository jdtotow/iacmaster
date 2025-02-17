package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CloudCredential struct {
	gorm.Model
	Name      string                 `gorm:"uniqueIndex"`
	Type      DestinationCloud       `json:"destination_cloud"`
	Variables map[string]string      `json:"variables" gorm:"type:jsonb"`
	ID        uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Settings  []IaCExecutionSettings `json:"settings"`
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
