package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `json:"name"`
	Uuid string `gorm:"primaryKey" json:"uuid"`
}

func (r Role) GetName() string {
	return r.Name
}
func (r *Role) SetName(name string) {
	r.Name = name
}

func (r *Role) SetUuid(uuid string) {
	r.Uuid = uuid
}
