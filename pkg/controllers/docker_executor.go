package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type DockerContainerController struct {
	cli                *client.Client
	ExecutorWorkingDir string
	Ports              []int
	SelectedPort       int
}

func NewDockerContainerController(working_dir string) *DockerContainerController {
	_client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &DockerContainerController{
		cli:                _client,
		ExecutorWorkingDir: working_dir,
		Ports:              []int{10000, 12000},
		SelectedPort:       0,
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
	port_str := os.Getenv("RUNNER_HOST_PORT")
	resp, err := d.cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
			Env:   env,
			Labels: map[string]string{
				"iacmaster": "true",
				"type":      "iacmaster/runner",
			},
			ExposedPorts: nat.PortSet{
				"8787/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			Mounts: volumeMounts,
			PortBindings: nat.PortMap{
				"8787/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: port_str,
					},
				},
			},
		},
		&network.NetworkingConfig{},
		nil,
		name,
	)
	if err != nil {
		return "", err
	}
	if err := d.cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		log.Fatalf("Error starting container: %v", err)
		return "", err
	}

	return resp.ID, nil
}

// Stop a container by name
func (d *DockerContainerController) StopContainerByID(containerID string) error {
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

func (d *DockerContainerController) GetExecutorHostIP() string {
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		return ""
	}
	var private_ip string
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsPrivate() {
				private_ip = ip.String()
			}
			// process IP address
		}
	}
	return private_ip
}

func (d *DockerContainerController) AddDeployment(deployment *msg.Deployment) (models.Executor, error) {
	volumes := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: d.ExecutorWorkingDir,
			Target: "/runner",
		},
	}

	system_address := os.Getenv("IACMASTER_SYSTEM_ADDRESS")
	executor := models.Executor{
		Name: deployment.EnvironmentID,
		State: models.ExecutorState{
			Status: models.InitStatus,
			Error:  nil,
		},
		Kind:             "docker",
		DepoymentID:      deployment.EnvironmentID,
		DeploymentObject: deployment,
	}
	if system_address == "" {
		log.Println("Please set IACMASTER_SYSTEM_ADDRESS")
		executor.State.Status = models.FailedStatus
		err := errors.New("master system address not set")
		executor.State.Error = err
		return executor, err
	}
	d.SelectedPort = rand.IntN(d.Ports[1]-d.Ports[0]) + d.Ports[0]
	var env_vars []string
	env_vars = append(env_vars, "IACMASTER_SYSTEM_ADDRESS="+system_address)
	env_vars = append(env_vars, "IACMASTER_SYSTEM_PORT="+os.Getenv("IACMASTER_SYSTEM_PORT"))
	env_vars = append(env_vars, "DEPLOYMENT_ID="+deployment.EnvironmentID)
	env_vars = append(env_vars, "EXECUTOR_HOST_IP="+d.GetExecutorHostIP())
	env_vars = append(env_vars, "RUNNER_HOST_PORT="+fmt.Sprintf("%v", d.SelectedPort))

	for key, value := range deployment.EnvironmentParameters {
		env_vars = append(env_vars, key+"="+value)
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
