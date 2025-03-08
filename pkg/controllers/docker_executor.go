package controllers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/jdtotow/iacmaster/pkg/models"
)

type DockerContainerController struct {
	cli                *client.Client
	ExecutorWorkingDir string
}

func NewDockerContainerController(working_dir string) *DockerContainerController {
	_client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &DockerContainerController{
		cli:                _client,
		ExecutorWorkingDir: working_dir,
	}
}

// List all running containers
func (d *DockerContainerController) ListRunningContainers(cli *client.Client) ([]container.Summary, error) {
	return d.cli.ContainerList(context.Background(), container.ListOptions{})
}

// Get container by name
func (d *DockerContainerController) GetContainerByID(containerID string) (*container.Summary, error) {
	containers, err := d.ListRunningContainers(d.cli)
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		if container.ID == containerID {
			return &container, nil
		}
	}

	return nil, nil
}

// Create a container
func (d *DockerContainerController) CreateContainer(name string, image string, env []string, volumeMounts []mount.Mount) (string, error) {
	resp, err := d.cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
			Env:   env,
		},
		&container.HostConfig{
			Mounts: volumeMounts,
		},
		&network.NetworkingConfig{},
		nil,
		name,
	)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// Stop a container by name
func (d *DockerContainerController) StopContainerByName(containerID string) error {
	_container, err := d.GetContainerByID(containerID)
	if err != nil {
		return err
	}
	if _container == nil {
		return nil
	}
	return d.cli.ContainerStop(context.Background(), _container.ID, container.StopOptions{})
}

// Remove a container by name
func (d *DockerContainerController) RemoveContainerByID(containerID string) error {
	_container, err := d.GetContainerByID(containerID)
	if err != nil {
		return err
	}
	if _container == nil {
		return nil
	}
	return d.cli.ContainerRemove(context.Background(), _container.ID, container.RemoveOptions{Force: true})
}

// Get container status by name
func (d *DockerContainerController) GetContainerStatus(containerID string) (string, error) {
	container, err := d.GetContainerByID(containerID)
	if err != nil {
		return "", err
	}
	if container == nil {
		return "not found", nil
	}
	return container.State, nil
}

func (d *DockerContainerController) AddDeployment(deployment models.Deployment) (models.Executor, error) {
	volumes := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: d.ExecutorWorkingDir,
			Target: "/tmp",
		},
	}
	var env_vars []string
	for key, value := range deployment.EnvironmentParameters {
		env_vars = append(env_vars, key+"="+value)
	}
	executor := models.Executor{
		Name: deployment.EnvironmentID,
		State: models.ExecutorState{
			Status: models.InitStatus,
			Error:  nil,
		},
		Kind:        "docker",
		DepoymentID: deployment.EnvironmentID,
	}

	containerID, err := d.CreateContainer(
		deployment.EnvironmentID,
		"iacmaster_runner:latest",
		env_vars,
		volumes,
	)
	if err != nil {
		executor.State.Status = models.FailedStatus
		executor.State.Error = err
		return executor, err
	}
	executor.ObjectID = containerID
	executor.State.Status = models.RunningStatus
	return executor, nil
}
func (d *DockerContainerController) RemoveDeployment(deploymentID string) error {
	return d.RemoveContainerByID(deploymentID)
}
