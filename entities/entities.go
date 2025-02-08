package entities

// User structure
type User struct {
	Fullname string
	Email    string
	Username string
	Password string
}

// Creating user
func CreateUser(fullname, email, username, password string) User {
	return User{
		Fullname: fullname,
		Email:    email,
		Username: username,
		Password: password,
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
}

func CreateUserGroup(name string) UserGroup {
	return UserGroup{Name: name, Members: []*User{}}
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
	Name string
}

func CreateOrganization(name string) Organization {
	return Organization{Name: name}
}

// Project structure
type Project struct {
	Name   string
	Parent string
	Org    Organization
}

func CreateProject(name, parent string, org Organization) Project {
	return Project{Name: name, Parent: parent, Org: org}
}
func (project Project) GetName() string {
	return project.Name
}
func (project Project) GetParent() string {
	return project.Parent
}
func (project Project) GetOrganization() Organization {
	return project.Org
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
