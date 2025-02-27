package controllers

import (
	"os"

	"github.com/jdtotow/iacmaster/pkg/models"
)

type Logic struct {
	Deployments        []*models.Deployment
	artifactController *IaCArtifactController
}

func CreateLogic(workingDir string) *Logic {
	return &Logic{
		Deployments:        []*models.Deployment{},
		artifactController: CreateIaCArtifactController(workingDir),
	}
}

func (l *Logic) AddDeployment(deployment *models.Deployment) bool {
	if l.HasDeployment(deployment.Name) {
		return false
	}
	l.Deployments = append(l.Deployments, deployment)
	err := l.GetRepo(*deployment)
	return err == nil
}

func (l *Logic) GetDeployments() []*models.Deployment {
	return l.Deployments
}

func (l *Logic) HasDeployment(name string) bool {
	for _, deployment := range l.Deployments {
		if deployment.Name == name {
			return true
		}
	}
	return false
}

func (l *Logic) GetRepo(deployment models.Deployment) error {
	localPath := l.artifactController.TmpFolderPath + "/" + deployment.GitData.Environment
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return l.artifactController.GetRepo(
			deployment.GitData.Url,
			deployment.GitData.Token,
			deployment.GitData.TokenUsername,
			deployment.GitData.Revision,
			deployment.GitData.ProxyUrl,
			deployment.GitData.ProxyUsername,
			deployment.GitData.ProxyPassword,
			deployment.GitData.Environment,
		)
	} else {
		return l.artifactController.UpdateRepo(
			deployment.GitData.Url,
			deployment.GitData.Token,
			deployment.GitData.TokenUsername,
			deployment.GitData.Revision,
			deployment.GitData.ProxyUrl,
			deployment.GitData.ProxyUsername,
			deployment.GitData.ProxyPassword,
			deployment.GitData.Environment,
		)
	}

}
