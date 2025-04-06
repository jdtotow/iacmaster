package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name      string            `gorm:"uniqueIndex" json:"name"`
	ID        uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Variables map[string]string `json:"variables" gorm:"type:jsonb"`
	Users     []User            `json:"users" gorm:"many2many:user_organizations;"`
}

func CreateOrganization(name string) Organization {
	return Organization{
		Name: name,
	}
}
func (org *Organization) SetName(name string) {
	org.Name = name
}
func (org *Organization) GetName() string {
	return org.Name
}
func (org *Organization) GetUsers() []User {
	return org.Users
}
func (org *Organization) HasUser(user User) bool {
	for _, _user := range org.Users {
		if _user.Username == user.Username {
			return true
		}
	}
	return false
}
func (org *Organization) AddUser(user User) {
	if !org.HasUser(user) {
		org.Users = append(org.Users, user)
	}
}
func (org *Organization) RemoveUser(user User) {
	for i, _user := range org.Users {
		if _user.Username == user.Username {
			org.Users = append(org.Users[:i], org.Users[i+1:]...)
			break
		}
	}
}
