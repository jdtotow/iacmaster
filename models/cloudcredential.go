package models

import "gorm.io/gorm"

type CloudCredential struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex"`
	Type      DestinationCloud
	Variables []EnvironmentVariable
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
