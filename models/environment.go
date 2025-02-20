package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnvStatus string

const EnvInit EnvStatus = "init"
const EnvUnknown EnvStatus = "unknown"
const EnvReady EnvStatus = "ready"
const EnvRunning EnvStatus = "running"
const EnvStopped EnvStatus = "stopped"
const EnvError EnvStatus = "error"
const EnvPending EnvStatus = "pending"

type Environment struct {
	gorm.Model
	Name                   string  `json:"name"`
	Project                Project `json:"project"`
	ProjectID              uuid.UUID
	IaCArtifact            IaCArtifact `json:"iac_artifact"`
	IaCArtifactID          uuid.UUID
	ExecSettings           IaCExecutionSettings
	IaCExecutionSettingsID uuid.UUID `json:"iac_execution_settings_uuid"`
	Status                 EnvStatus `json:"status"`
	ID                     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

func (env *Environment) SetName(name string) {
	env.Name = name
}
func (env Environment) GetName() string {
	return env.Name
}
