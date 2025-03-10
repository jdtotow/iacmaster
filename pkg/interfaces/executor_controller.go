package interfaces

import "github.com/jdtotow/iacmaster/pkg/models"

type ExecutorController interface {
	AddDeployment(deployment models.Deployment) (models.Executor, error)
	RemoveDeployment(deploymentID string) error
}
