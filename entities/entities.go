package entities

import (
	"github.com/google/uuid"
)

// User structure
type User struct {
	Fullname       string
	Email          string
	Username       string
	Password       string
	OrganizationId string
	ID             string
}

// Creating user
func CreateUser(fullname, email, username, password string) User {
	return User{
		Fullname: fullname,
		Email:    email,
		Username: username,
		Password: password,
		ID:       uuid.New().String(),
	}
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

// ///////////////////////
// User Group structure
type UserGroup struct {
	Name    string
	Members []*User
	ID      string
}

func CreateUserGroup(name string) UserGroup {
	return UserGroup{
		Name:    name,
		Members: []*User{},
		ID:      uuid.New().String(),
	}
}
func (group UserGroup) GetName() string {
	return group.Name
}
func (group UserGroup) GetMembers() []*User {
	return group.Members
}
func (group *UserGroup) SetName(name string) {
	group.Name = name
}
func (group UserGroup) HasMember(user User) bool {
	for _, _user := range group.Members {
		if _user.GetUsername() == user.GetUsername() {
			return true
		}
	}
	return false
}
func (group *UserGroup) AddMember(user *User) {
	if !group.HasMember(*user) {
		group.Members = append(group.Members, user)
	}
}

// Organization structure
type Organization struct {
	Name      string
	Admin     *User
	ID        string
	Variables []*EnvironmentVariable
}

func CreateOrganization(name string) Organization {
	return Organization{
		Name: name,
		ID:   uuid.New().String(),
	}
}
func (org *Organization) SetName(name string) {
	org.Name = name
}
func (org *Organization) SetAdmin(user *User) {
	org.Admin = user
	user.OrganizationId = org.ID
}
func (org Organization) GetOrgID() string {
	return org.ID
}

// Project structure
type Project struct {
	Name           string
	Parent         string
	OrganizationID string
	ID             string
	Variables      []*EnvironmentVariable
}

func CreateProject(name, parent string, org string) Project {
	return Project{
		Name:           name,
		Parent:         parent,
		OrganizationID: org,
		ID:             uuid.New().String(),
	}
}
func (project Project) GetName() string {
	return project.Name
}
func (project Project) GetParent() string {
	return project.Parent
}
func (project Project) GetOrganization() string {
	return project.OrganizationID
}

// Environment variables
type EnvironmentVariable struct {
	Name  string
	Value interface{}
}

func CreateEnvironmentVariable(name string, value interface{}) EnvironmentVariable {
	return EnvironmentVariable{
		Name:  name,
		Value: value,
	}
}
func (env EnvironmentVariable) GetName() string {
	return env.Name
}
func (env EnvironmentVariable) GetValue() interface{} {
	return env.Value
}
func (env *EnvironmentVariable) SetName(name string) {
	env.Name = name
}
func (env *EnvironmentVariable) SetValue(value interface{}) {
	env.Value = value
}

// IacArtifact
type IaCArtifact struct {
	Type   string
	Name   string
	ScmUrl string
	ID     string
}

func CreateIaCArtifact(_type, name, scm string) IaCArtifact {
	return IaCArtifact{
		Type:   _type,
		Name:   name,
		ScmUrl: scm,
	}
}
func (arti IaCArtifact) GetType() string {
	return arti.Type
}
func (arti IaCArtifact) GetName() string {
	return arti.Name
}
func (arti IaCArtifact) GetSCMUrl() string {
	return arti.ScmUrl
}
func (arti *IaCArtifact) SetType(_type string) {
	arti.Type = _type
}
func (arti *IaCArtifact) SetName(name string) {
	arti.Name = name
}
func (arti *IaCArtifact) SetSCMurl(url string) {
	arti.ScmUrl = url
}

// Environment status structure
type EnvironmentStatus string

// Environment structure
const (
	Inactive EnvironmentStatus = "inactive"
	Active   EnvironmentStatus = "active"
	Running  EnvironmentStatus = "running"
	Stopped  EnvironmentStatus = "stopped"
)

type Environment struct {
	Name          string
	ProjectID     string
	IaCArtifactID string
	Status        EnvironmentStatus
}
