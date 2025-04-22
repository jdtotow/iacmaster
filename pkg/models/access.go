package models

import "gorm.io/gorm"

type ResourceType string

const (
	OrganizationResourceType    ResourceType = "organization"
	GroupResourceType           ResourceType = "group"
	ProjectResourceType         ResourceType = "project"
	UserResourceType            ResourceType = "user"
	RoleResourceType            ResourceType = "role"
	CloudCredentialResourceType ResourceType = "cloud_credential"
	IaCArtifactResourceType     ResourceType = "iac_artifact"
)

type Access struct {
	gorm.Model
	UserID         string       `json:"user_id"`
	RoleID         string       `json:"role_id"`
	OrganizationID string       `json:"organization_id"`
	GroupID        string       `json:"group_id"`
	ResourceType   ResourceType `json:"resource_type"`
	ResourceID     string       `json:"resource_id"`
}

func NewAccess(userID, roleID, organizationID, groupID, resourceID string, resourceType ResourceType) *Access {
	return &Access{
		UserID:         userID,
		RoleID:         roleID,
		OrganizationID: organizationID,
		GroupID:        groupID,
		ResourceType:   resourceType,
		ResourceID:     resourceID,
	}
}

func (a *Access) SetUserID(user User) {
	a.UserID = user.ID.String()
}

func (a *Access) SetRoleID(role Role) {
	a.RoleID = role.ID.String()
}

func (a *Access) SetOrganizationID(organization Organization) {
	a.OrganizationID = organization.ID.String()
}
func (a *Access) SetGroupID(group Group) {
	a.GroupID = group.ID.String()
}
func (a *Access) SetResourceID(resourceID string) {
	a.ResourceID = resourceID
}
func (a *Access) SetResourceType(resourceType string) {
	a.ResourceType = ResourceType(resourceType)
}
