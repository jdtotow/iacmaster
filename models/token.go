package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenType string

const GitToken TokenType = "git"
const AuthorizationToken TokenType = "authorization"

type Token struct {
	gorm.Model
	ID                   uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name                 string                 `json:"name"`
	Type                 TokenType              `json:"type"`
	Token                string                 `json:"token"`
	Project              Project                `json:"project"`
	ProjectID            uuid.UUID              `json:"project_id"`
	IaCExecutionSettings []IaCExecutionSettings `json:"settings"`
}
