package worker

type KubernetesRunner struct {
	Name   string
	Type   string
	Status JobStatus
}

func CreateKubernetesRunner(name string) *KubernetesRunner {
	return &KubernetesRunner{
		Name:   name,
		Type:   "kubernetes",
		Status: "init",
	}
}

func (d KubernetesRunner) GetName() string {
	return d.Name
}

func (d KubernetesRunner) GetType() string {
	return d.Type
}

func (d KubernetesRunner) GetJobStatus() string {
	return string(d.Status)
}

func (d *KubernetesRunner) StartJob() error {
	return nil
}

func (d *KubernetesRunner) StopJob() error {
	return nil
}
