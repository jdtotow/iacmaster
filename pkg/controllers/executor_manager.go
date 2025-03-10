package controllers

import (
	"fmt"
	"log"

	"github.com/jdtotow/iacmaster/pkg/interfaces"
	"github.com/jdtotow/iacmaster/pkg/models"
)

type ExecutorManager struct {
	executors          []*models.Executor
	MaxExecutors       int
	Available          int
	ExecutorWorkingDir string
	ExecutionPlatform  string
}

func CreateExecutorManager(working_dir, ExecutionPlatform string) *ExecutorManager {
	return &ExecutorManager{
		executors:          []*models.Executor{},
		MaxExecutors:       10,
		ExecutorWorkingDir: working_dir,
		ExecutionPlatform:  ExecutionPlatform,
	}
}

func (em *ExecutorManager) CreateExecutorController() interfaces.ExecutorController {
	switch em.ExecutionPlatform {
	case "docker":
		return NewDockerContainerController(em.ExecutorWorkingDir)
	case "kubernetes":
		return NewKubernetesPodController("iacmaster", em.ExecutorWorkingDir)
	default:
		return nil
	}
}

func (em *ExecutorManager) AddExecutor(executor *models.Executor) error {
	if em.ExecutorExists(executor.Name) {
		return fmt.Errorf("executor already exists")
	}
	if len(em.executors) >= em.MaxExecutors {
		return fmt.Errorf("max executors reached")
	}
	em.executors = append(em.executors, executor)
	return nil
}
func (em *ExecutorManager) RemoveExecutor(executor *models.Executor) {
	for i, e := range em.executors {
		if e == executor {
			em.executors = append(em.executors[:i], em.executors[i+1:]...)
			break
		}
	}
}
func (em *ExecutorManager) GetExecutor(name string) *models.Executor {
	for _, e := range em.executors {
		if e.Name == name {
			return e
		}
	}
	return nil
}
func (em *ExecutorManager) GetExecutors() []*models.Executor {
	return em.executors
}
func (em *ExecutorManager) GetMaxExecutors() int {
	return em.MaxExecutors
}
func (em *ExecutorManager) ExecutorExists(name string) bool {
	for _, e := range em.executors {
		if e.Name == name {
			return true
		}
	}
	return false
}

func (em *ExecutorManager) StartDeployment(deployment models.Deployment) error {
	executorController := em.CreateExecutorController()
	if executorController == nil {
		log.Println(em.ExecutionPlatform)
		return fmt.Errorf("execution platform not supported")
	}
	executor, err := executorController.AddDeployment(deployment)
	if err != nil {
		return err
	}
	err = em.AddExecutor(&executor)
	if err != nil {
		return err
	}
	return nil
}
func (em *ExecutorManager) SetExecutorState(executor_name string, state models.ExecutorState) {
	executor := em.GetExecutor(executor_name)
	if executor == nil {
		return
	}
	executor.State = state
}
