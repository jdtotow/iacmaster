package controllers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DockerContainerController struct {
	cli *client.Client
}

func NewDockerContainerController() *DockerContainerController {
	_client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &DockerContainerController{
		cli: _client,
	}
}

// List all running containers
func (d *DockerContainerController) ListRunningContainers(cli *client.Client) ([]container.Summary, error) {
	return d.cli.ContainerList(context.Background(), container.ListOptions{})
}

// Get container by name
func (d *DockerContainerController) GetContainerByName(name string) (*container.Summary, error) {
	containers, err := d.ListRunningContainers(d.cli)
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == "/"+name {
				return &container, nil
			}
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
func (d *DockerContainerController) StopContainerByName(name string) error {
	_container, err := d.GetContainerByName(name)
	if err != nil {
		return err
	}
	if _container == nil {
		return nil
	}
	return d.cli.ContainerStop(context.Background(), _container.ID, container.StopOptions{})
}

// Remove a container by name
func (d *DockerContainerController) RemoveContainerByName(name string) error {
	_container, err := d.GetContainerByName(name)
	if err != nil {
		return err
	}
	if _container == nil {
		return nil
	}
	return d.cli.ContainerRemove(context.Background(), _container.ID, container.RemoveOptions{Force: true})
}

// Get container status by name
func (d *DockerContainerController) GetContainerStatus(name string) (string, error) {
	container, err := d.GetContainerByName(name)
	if err != nil {
		return "", err
	}
	if container == nil {
		return "not found", nil
	}
	return container.State, nil
}
