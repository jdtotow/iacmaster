package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name   string
	UserID uint
}

func (r Role) GetName() string {
	return r.Name
}
func (r *Role) SetName(name string) {
	r.Name = name
}
