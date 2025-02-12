package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Fullname       string
	Email          string `gorm:"uniqueIndex"`
	Username       string
	Password       string
	OrganizationId string
	groups         []UserGroup
	roles          []Role
}

// GetFullname
func (user User) GetFullname() string {
	return user.Fullname
}

// Get user email
func (user User) GetEmail() string {
	return user.Email
}
func (user User) GetUsername() string {
	return user.Username
}
func (user User) GetPassword() string {
	return user.Password
}
func (user *User) SetFullname(fullname string) {
	user.Fullname = fullname
}
func (user *User) SetEmail(email string) {
	user.Email = email
}
func (user *User) SetUsername(username string) {
	user.Username = username
}
func (user *User) SetPassword(password string) {
	user.Password = password
}
func (user *User) AssignToGroup(group UserGroup) {
	user.groups = append(user.groups, group)
}
func (user User) IsMemberOfGroup(group UserGroup) bool {
	for _, _group := range user.groups {
		if _group.Name == group.Name {
			return true
		}
	}
	return false
}
func (user User) HasRole(role Role) bool {
	for _, _role := range user.roles {
		if _role.Name == role.Name {
			return true
		}
	}
	return false
}
