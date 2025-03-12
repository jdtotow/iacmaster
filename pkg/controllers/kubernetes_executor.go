package controllers

import (
	"context"
	"fmt"

	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesPodController struct {
	clientset          *kubernetes.Clientset
	namespace          string
	executorWotkingDir string
}

func NewKubernetesPodController(namespace, executor_working_dir string) *KubernetesPodController {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil
	}
	return &KubernetesPodController{
		clientset:          clientset,
		namespace:          namespace,
		executorWotkingDir: executor_working_dir,
	}
}

func (k *KubernetesPodController) CreatePod(namespace, name, image string, envVars map[string]string, volumes []v1.Volume) (string, error) {
	env := []v1.EnvVar{}
	for key, value := range envVars {
		env = append(env, v1.EnvVar{Name: key, Value: value})
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  name,
					Image: image,
					Env:   env,
				},
			},
			Volumes: volumes,
		},
	}

	_, err := k.clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	fmt.Printf("Pod %s created successfully in namespace %s\n", name, namespace)
	return pod.Name, nil
}

func (k *KubernetesPodController) GetPod(namespace, podName string) (*v1.Pod, error) {
	return k.clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
}

func (k *KubernetesPodController) GetAllPods(namespace string) ([]v1.Pod, error) {
	podList, err := k.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

func (k *KubernetesPodController) DeletePod(namespace, podName string) error {
	return k.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}

func (d *KubernetesPodController) AddDeployment(deployment *msg.Deployment) (models.Executor, error) {

	volumes := v1.Volume{
		Name: "workspace",
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: d.executorWotkingDir,
			},
		},
	}

	executor := models.Executor{
		Name:             deployment.EnvironmentID,
		Status:           models.InitStatus,
		Kind:             "kubernetes",
		DepoymentID:      deployment.EnvironmentID,
		DeploymentObject: deployment,
	}
	docker_image := "iacmaster_runner:latest"
	pod_name, err := d.CreatePod(d.namespace, deployment.EnvironmentID, docker_image, deployment.EnvironmentParameters, []v1.Volume{volumes})
	if err != nil {
		executor.Status = models.FailedStatus
		executor.Error = err.Error()
		return executor, err
	}
	executor.ObjectID = pod_name
	executor.Status = models.RunningStatus
	return executor, nil
}
func (d *KubernetesPodController) RemoveDeployment(podName string) error {
	return d.DeletePod(d.namespace, podName)
}
