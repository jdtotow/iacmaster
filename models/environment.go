package models

import "gorm.io/gorm"

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
	Name                     string               `json:"name"`
	Project                  Project              `gorm:"foreignKey:Uuid;references:ProjectUuid" json:"project"`
	ProjectUuid              string               `json:"projectuuid"`
	Artifact                 IaCArtifact          `gorm:"foreignKey:Uuid;references:IaCArtifactUuid" json:"artifact"`
	IaCArtifactUuid          string               `json:"iacartifact_uuid"`
	ExecSettings             IaCExecutionSettings `gorm:"foreignKey:Uuid;references:IaCExecutionSettingsUuid"`
	IaCExecutionSettingsUuid string               `json:"iac_execution_settings_uuid"`
	Status                   EnvStatus            `json:"status"`
	Uuid                     string               `gorm:"primaryKey" json:"uuid"`
}

func (env *Environment) SetName(name string) {
	env.Name = name
}
func (env Environment) GetName() string {
	return env.Name
}
func (env *Environment) SetProject(project Project) {
	env.Project = project
}
func (env *Environment) SetArtifact(arti IaCArtifact) {
	env.Artifact = arti
}
func (env *Environment) SetExecutionSettings(execSettings IaCExecutionSettings) {
	env.ExecSettings = execSettings
}
func (env *Environment) SetUuid(uuid string) {
	env.Uuid = uuid
}
