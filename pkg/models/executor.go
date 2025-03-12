package models

import "github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"

type ExecutorKind string

const ShellExecutor ExecutorKind = "shell"
const DockerExecutor ExecutorKind = "docker"
const KubernetesExecutor ExecutorKind = "kubernetes"

type ExecutorStatus string

const InitStatus ExecutorStatus = "init"
const RunningStatus ExecutorStatus = "running"
const FailedStatus ExecutorStatus = "failed"
const SucceededStatus ExecutorStatus = "succeeded"

type ExecutorController interface {
	AddExecutor(executor Executor) error
	RemoveExecutor(executor Executor)
	GetExecutor(name string) Executor
	GetExecutors() []Executor
	GetMaxExecutors() int
	ExecutorExists(name string) bool
	StartExecutor(deployment msg.Deployment) (Executor, error)
}

type Executor struct {
	Kind             ExecutorKind
	Status           ExecutorStatus
	Name             string
	DepoymentID      string
	ObjectID         string //the object can be a container, a pod or a process pid
	DeploymentObject *msg.Deployment
	Error            string
}

func (e *Executor) SetError(err error) {
	e.Error = err.Error()
}
func (e *Executor) SetStatus(status ExecutorStatus) {
	e.Status = status
}
func (e *Executor) SetName(name string) {
	e.Name = name
}

func (e *Executor) SetDeploymentID(id string) {
	e.DepoymentID = id
}
func (e *Executor) GetKind() ExecutorKind {
	return e.Kind
}
func (e *Executor) GetName() string {
	return e.Name
}
func (e *Executor) GetDeploymentID() string {
	return e.DepoymentID
}
func (e *Executor) GetState() ExecutorStatus {
	return e.Status
}
