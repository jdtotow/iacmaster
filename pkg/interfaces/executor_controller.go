package interfaces

import (
	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type ExecutorController interface {
	AddDeployment(deployment *msg.Deployment) (models.Executor, error)
	RemoveDeployment(deploymentID string) error
}
