package controllers

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"reflect"

	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type IaCRunner struct {
	Deployment         *msg.Deployment
	artifactController *IaCArtifactController
	Name               string
	Kind               models.ExecutorKind
	State              models.ExecutorState
	EnvironmentID      string
	MandatoryCommands  []string
}

func CreateIaCRunner(workingDir, name string, mandatory_commands []string, kind models.ExecutorKind) *IaCRunner {
	return &IaCRunner{
		Name:              name,
		Kind:              kind,
		MandatoryCommands: mandatory_commands,
		State: models.ExecutorState{
			Status: models.InitStatus,
			Error:  nil,
		},
		artifactController: CreateIaCArtifactController(workingDir),
	}
}

func (l *IaCRunner) DeleteDeployment(deployment *msg.Deployment) {
	localPath := l.artifactController.TmpFolderPath + "/" + deployment.EnvironmentID + "/" + deployment.HomeFolder
	err := l.terraformDestroy(localPath)
	if err != nil {
		log.Println(err.Error())
	}
	os.RemoveAll(localPath)
}
func (l *IaCRunner) SetDeployment(deployment *msg.Deployment) bool {
	l.Deployment = deployment
	if deployment.TerraformVersion != "" {
		err := l.setTerraformVersion(deployment.TerraformVersion)
		if err != nil {
			log.Println(err.Error())
			deployment.Error = err.Error()
			l.State.Status = models.FailedStatus
			l.State.Error = err
			return false
		}
	}
	l.Deployment = deployment
	err := l.GetRepo(deployment)
	if err != nil {
		log.Println("Error -> ", err.Error())
		deployment.Error = err.Error()
		l.State.Status = models.FailedStatus
		l.State.Error = err
		return false
	}
	//
	if deployment.CloudDestination == "azure" {
		err = l.azureLogin(deployment)
		if err != nil {
			log.Println("An error occured on deployment ", deployment.EnvironmentID, " :", err.Error())
			deployment.Error = err.Error()
			return false
		} else {
			log.Println("Deployment ", deployment.EnvironmentID, " login succeeded on azure")
			//deployment.AddActivity("Azure login succeeded")
		}
	} else if deployment.CloudDestination == "aws" {
		os.Setenv("AWS_ACCESS_KEY_ID", deployment.EnvironmentParameters["AWS_ACCESS_KEY_ID"])
		os.Setenv("AWS_SECRET_ACCESS_KEY", deployment.EnvironmentParameters["AWS_SECRET_ACCESS_KEY"])
		log.Println("Deployment ", deployment.EnvironmentID, " AWS access parameters define")
	} else if deployment.CloudDestination == "gcp" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", deployment.EnvironmentParameters["GOOGLE_APPLICATION_CREDENTIALS"])
		log.Println("Deployment ", deployment.EnvironmentID, " GCP credential file defined")
	} else {
		deployment.Error = "Cloud : " + deployment.CloudDestination + " is not supported"
		log.Println("Deployment ", deployment.EnvironmentID, " , Cloud : "+deployment.CloudDestination+" is not supported")
		return false
	}
	err = l.deployEnvironment(deployment)
	if err != nil {
		deployment.Error = err.Error()
		log.Println("Deployment ", deployment.EnvironmentID, " deployment failed : ", err.Error())
		return false
	} else {
		log.Println("Deployment ", deployment.EnvironmentID, " deployment succeeded")
	}
	l.State.Status = models.SucceededStatus
	l.State.Error = nil
	return true
}

func (l *IaCRunner) deployEnvironment(deployment *msg.Deployment) error {
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

func (l *IaCRunner) GetDeployment() *msg.Deployment {
	return l.Deployment
}

func (l *IaCRunner) GetRepo(deployment *msg.Deployment) error {
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

func (l *IaCRunner) CheckIfMandatoryCommandExists(commands string) bool {
	for _, command := range l.MandatoryCommands {
		_, err := exec.LookPath(command)
		if err != nil {
			log.Println("Command not found : ", command)
			l.State.Status = models.FailedStatus
			l.State.Error = err
			return false
		}
	}
	return true
}

func (l *IaCRunner) runCommand(prog string, commands []string) error {
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

func (l *IaCRunner) setTerraformVersion(terraform_version string) error {
	prog := "tfenv"
	var commands []string
	commands = append(commands, "use")
	commands = append(commands, terraform_version)
	err := l.runCommand(prog, commands)
	return err
}

func (l *IaCRunner) azureLogin(deployment *msg.Deployment) error {
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
	if err != nil {
		log.Println(err.Error())
	}

	commands = commands[:0]
	commands = append(commands, "account")
	commands = append(commands, "set")
	commands = append(commands, "--subscription")
	commands = append(commands, deployment.EnvironmentParameters["ARM_SUBSCRIPTION_ID"])

	err = l.runCommand(prog, commands)
	return err
}

func (l *IaCRunner) terraformInit(folder string) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "init")
	err := l.runCommand(prog, commands)
	return err
}

func (l *IaCRunner) terraformPlan(folder, var_file_path string, save bool) error {
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

func (l *IaCRunner) terraformApply(folder, var_file_path string, saved bool) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "apply")
	commands = append(commands, "-auto-approve")

	if saved {
		commands = append(commands, "plan.tfplan")
	}
	err := l.runCommand(prog, commands)
	return err
}

func (l *IaCRunner) terraformDestroy(folder string) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "destroy")
	commands = append(commands, "-auto-approve")
	err := l.runCommand(prog, commands)
	return err
}

func (s *IaCRunner) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		log.Println("Runner actor started on address -> ", ctx.Engine().Address())
		systemPID := actor.NewPID("192.168.1.128:3434", "iacmaster/system")
		ctx.Send(systemPID, &msg.RunnerStatus{Name: s.Name, Status: "Ready", Address: "192.168.1.128:8787"})
	case actor.Initialized:
		log.Println("Runner actor initialized")
	case *actor.PID:
		log.Println("Runner actor has god an ID")
	case *msg.Deployment:
		log.Println("Depoyment object received")
		s.SetDeployment(m)
	default:
		slog.Warn("server got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}
