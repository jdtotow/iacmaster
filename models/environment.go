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
	ID                       uint
	Name                     string               `json:"name"`
	Project                  Project              `gorm:"foreignKey:ProjectUuid" json:"project"`
	ProjectUuid              string               `json:"projectuuid"`
	Artifact                 IaCArtifact          `gorm:"foreignKey:Uuid;references:IaCArtifactUuid" json:"artifact"`
	IaCArtifactUuid          uuid.UUID            `json:"iacartifact_uuid"`
	ExecSettings             IaCExecutionSettings `gorm:"foreignKey:IaCExecutionSettingsUuid"`
	IaCExecutionSettingsUuid uuid.UUID            `json:"iac_execution_settings_uuid"`
	Status                   EnvStatus            `json:"status"`
	Uuid                     uuid.UUID            `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

func (env *Environment) SetName(name string) {
	env.Name = name
}
func (env Environment) GetName() string {
	return env.Name
}
