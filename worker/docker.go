package worker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type DockerRunner struct {
	Name       string
	Type       string
	Status     JobStatus
	CurrentJob JobData
}

func CreateDockerRunner(name string) DockerRunner {
	return DockerRunner{
		Name:   name,
		Type:   "docker",
		Status: "init",
	}
}

func (d DockerRunner) GetName() string {
	return d.Name
}

func (d DockerRunner) GetType() string {
	return d.Type
}

func (d DockerRunner) GetJobStatus() JobStatus {
	return d.Status
}

func (d DockerRunner) StartJob() error {
	return nil
}

func (d DockerRunner) StopJob() error {
	return nil
}

func (d DockerRunner) SetJobInfo(data JobData) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer cli.Close()
	ctx := context.Background()
	parameters := []string{}
	parameters = append(parameters, "TERRAFORM_VERSION="+data.TerraformVersion)
	parameters = append(parameters, "WORKING_DIR="+data.WorkingDir)

	for name, value := range data.EnvironmentParameters {
		parameters = append(parameters, name+"="+value)
	}
	commands := []string{}
	commands = append(commands, "/app/entrypoint.sh")

	containerConfig := &container.Config{
		Image: data.DockerImage,
		Env:   parameters,
		Cmd:   commands,
	}
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{ // Define volume mounts
			{
				Type:   mount.TypeBind,
				Source: data.VolumePath, // Change this to a valid host path
				Target: data.WorkingDir, // Mount point inside container
			},
			{
				Type:   mount.TypeBind,
				Source: "/Users/jean-didier/Projects/IaCMaster/script",
				Target: "/app",
			},
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}
	/*
		_, err = cli.ImagePull(ctx, data.DockerImage, image.PullOptions{})
		if err != nil {
			log.Fatalf("Error pulling image: %v", err)
		}
		fmt.Println("Image pulled successfully!")
	*/

	// Create the container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "docker-worker-terraform")
	if err != nil {
		log.Fatalf("Error creating container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatalf("Error starting container: %v", err)
	}
}

func (d DockerRunner) GetWorkerInfo() WorkerInfo {
	return WorkerInfo{
		hostname: "localhost",
	}
}
