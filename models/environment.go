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
	Name                   string `json:"name"`
	Project                Project
	ProjectID              int
	Artifact               IaCArtifact
	IaCArtifactID          int
	ExecSettings           IaCExecutionSettings
	IaCExecutionSettingsID int
	Status                 EnvStatus
	Uuid                   string
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
