package controllers

import (
	"github.com/jdtotow/iacmaster/worker"
)

type Runner struct {
}

func (r *Runner) Start(name, _type string) worker.Worker {
	if _type == "kubernetes" {
		return worker.CreateKubernetesRunner(name)
	} else {
		return worker.CreateDockerRunner(name)
	}
}
