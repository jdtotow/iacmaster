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
	OrganizationID int
	Organization   Organization
	UserGroupID    int
	groups         []UserGroup
	RoleID         int
	roles          []Role
	Uuid           string
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
	if !user.IsMemberOfGroup(group) {
		user.groups = append(user.groups, group)
	}
}
func (user User) IsMemberOfGroup(group UserGroup) bool {
	for _, _group := range user.groups {
		if _group.Name == group.Name {
			return true
		}
	}
	return false
}
func (user *User) AddRole(role Role) {
	for _, _role := range user.roles {
		if role.GetName() == _role.GetName() {
			return
		}
	}
	user.roles = append(user.roles, role)
}
func (user User) HasRole(role Role) bool {
	for _, _role := range user.roles {
		if _role.Name == role.Name {
			return true
		}
	}
	return false
}
func (user *User) SetOrganization(org Organization) {
	user.Organization = org
}
func (user User) GetOrganization() Organization {
	return user.Organization
}
func (user *User) SetUuid(uuid string) {
	user.Uuid = uuid
}
