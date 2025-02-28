package controllers

import (
	"bufio"
	"log"
	"os"
	"os/exec"

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
	if !l.HasDeployment(deployment.Name) {
		l.Deployments = append(l.Deployments, deployment)
	}
	err := l.GetRepo(*deployment)
	if err != nil {
		deployment.SetError(err.Error())
		return false
	}
	deployment.AddActivity("Repository cloned or updated")
	if deployment.CloudDestination == "azure" {
		err = l.azureLogin(deployment)
		if err != nil {
			log.Println("An error occured on deployment ", deployment.EnvironmentID, " :", err.Error())
			deployment.SetError(err.Error())
			return false
		} else {
			log.Println("Deployment ", deployment.EnvironmentID, " login succeeded on azure")
			deployment.AddActivity("Azure login succeeded")
		}
	} else if deployment.CloudDestination == "aws" {
		os.Setenv("AWS_ACCESS_KEY_ID", deployment.EnvironmentParameters["AWS_ACCESS_KEY_ID"])
		os.Setenv("AWS_SECRET_ACCESS_KEY", deployment.EnvironmentParameters["AWS_SECRET_ACCESS_KEY"])
		log.Println("Deployment ", deployment.EnvironmentID, " AWS access parameters define")
		deployment.AddActivity("AWS access parameters defined")
	} else if deployment.CloudDestination == "gcp" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", deployment.EnvironmentParameters["GOOGLE_APPLICATION_CREDENTIALS"])
		deployment.AddActivity("GCP credential file defined")
		log.Println("Deployment ", deployment.EnvironmentID, " GCP credential file defined")
	} else {
		deployment.SetError("Cloud : " + deployment.CloudDestination + " is not supported")
		log.Println("Deployment ", deployment.EnvironmentID, " , Cloud : "+deployment.CloudDestination+" is not supported")
		return false
	}
	err = l.deployEnvironment(deployment)
	if err != nil {
		deployment.SetError(err.Error())
		log.Println("Deployment ", deployment.EnvironmentID, " deployment failed : ", err.Error())
		return false
	} else {
		log.Println("Deployment ", deployment.EnvironmentID, " deployment succeeded")
	}
	return true
}

func (l *Logic) deployEnvironment(deployment *models.Deployment) error {
	localPath := l.artifactController.TmpFolderPath + "/" + deployment.EnvironmentID + "/" + deployment.HomeFolder
	log.Println(localPath)
	err := l.terraformInit(localPath)
	if err != nil {
		return err
	}
	err = l.terraformPlan(localPath, "", true)
	if err != nil {
		return err
	}
	err = l.terraformApply(localPath, "", true)
	if err != nil {
		return err
	}
	return nil
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
	localPath := l.artifactController.TmpFolderPath + "/" + deployment.EnvironmentID
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return l.artifactController.GetRepo(
			deployment.GitData.Url,
			deployment.GitData.Token,
			deployment.GitData.TokenUsername,
			deployment.GitData.Revision,
			deployment.GitData.ProxyUrl,
			deployment.GitData.ProxyUsername,
			deployment.GitData.ProxyPassword,
			deployment.EnvironmentID,
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
			deployment.EnvironmentID,
		)
	}

}

func (l *Logic) runCommand(prog string, commands []string) error {
	cmd := exec.Command(prog, commands...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		log.Println(in.Text()) // write each line to your log, or anything you need
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	return nil
}

func (l *Logic) azureLogin(deployment *models.Deployment) error {
	//login
	//az login --service-principal --username $ARM_CLIENT_ID --password $ARM_CLIENT_SECRET --tenant $ARM_TENANT_ID
	log.Println("Azure login ...")
	prog := "az"
	var commands []string
	commands = append(commands, "login")
	commands = append(commands, "--service-principal")
	commands = append(commands, "--username")
	commands = append(commands, deployment.EnvironmentParameters["ARM_CLIENT_ID"])
	commands = append(commands, "--password")
	commands = append(commands, deployment.EnvironmentParameters["ARM_CLIENT_SECRET"])
	commands = append(commands, "--tenant")
	commands = append(commands, deployment.EnvironmentParameters["ARM_TENANT_ID"])

	err := l.runCommand(prog, commands)
	return err
}

func (l *Logic) terraformInit(folder string) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "init")
	err := l.runCommand(prog, commands)
	return err
}

func (l *Logic) terraformPlan(folder, var_file_path string, save bool) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "plan")
	if save {
		commands = append(commands, "-out=plan.tfplan")
	}
	if var_file_path != "" {
		commands = append(commands, "-var-file="+var_file_path)
	}
	err := l.runCommand(prog, commands)
	return err
}

func (l *Logic) terraformApply(folder, var_file_path string, saved bool) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "apply")
	if var_file_path != "" {
		commands = append(commands, "-var-file="+var_file_path)
	}
	if saved {
		commands = append(commands, "plan.tfplan")
	}
	err := l.runCommand(prog, commands)
	return err
}

func (l *Logic) terraformDestroy(folder string) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "destroy")
	err := l.runCommand(prog, commands)
	return err
}
