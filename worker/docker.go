package worker

type DockerRunner struct {
	Name   string
	Type   string
	Status JobStatus
}

func CreateDockerRunner(name string) *DockerRunner {
	return &DockerRunner{
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

func (d DockerRunner) GetJobStatus() string {
	return string(d.Status)
}

func (d *DockerRunner) StartJob() error {
	return nil
}

func (d *DockerRunner) StopJob() error {
	return nil
}
