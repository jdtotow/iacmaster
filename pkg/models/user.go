package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Fullname      string         `json:"fullname"`
	Email         string         `gorm:"uniqueIndex" json:"email"`
	Username      string         `gorm:"uniqueIndex" json:"username"`
	Password      string         `json:"password"`
	Organizations []Organization `json:"organization" gorm:"many2many:user_organizations;"`
	Groups        []Group        `json:"groups" gorm:"many2many:user_groups;"`
	Roles         []Role         `json:"roles" gorm:"many2many:user_roles;"`
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
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
func (user *User) AssignToGroup(group Group) {
	if !user.IsMemberOfGroup(group) {
		user.Groups = append(user.Groups, group)
	}
}
func (user User) IsMemberOfGroup(group Group) bool {
	for _, _group := range user.Groups {
		if _group.Name == group.Name {
			return true
		}
	}
	return false
}
func (user *User) AddRole(role Role) {
	for _, _role := range user.Roles {
		if role.GetName() == _role.GetName() {
			return
		}
	}
	user.Roles = append(user.Roles, role)
}
func (user User) HasRole(role Role) bool {
	for _, _role := range user.Roles {
		if _role.Name == role.Name {
			return true
		}
	}
	return false
}
func (user *User) IsMemberOfOrganization(org Organization) bool {
	for _, _org := range user.Organizations {
		if _org.Name == org.Name {
			return true
		}
	}
	return false
}
func (user *User) AddOrganization(org Organization) {
	user.Organizations = append(user.Organizations, org)
}
func (user User) GetOrganizations() []Organization {
	return user.Organizations
}
