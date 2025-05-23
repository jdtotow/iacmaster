package controllers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type IaCRunner struct {
	Deployment         *msg.Deployment
	artifactController *IaCArtifactController
	Name               string
	Kind               models.ExecutorKind
	Status             models.ExecutorStatus
	EnvironmentID      string
	MandatoryCommands  []string
	Engine             *actor.Engine
	CommandTimeout     int //in minute
	SystemPID          *actor.PID
	Error              error
}

func CreateIaCRunner(workingDir, name string, mandatory_commands []string, kind models.ExecutorKind, engine *actor.Engine) *IaCRunner {
	systemAddr := os.Getenv("IACMASTER_SYSTEM_ADDRESS") + ":" + os.Getenv("IACMASTER_SYSTEM_PORT")
	systemPID := actor.NewPID(systemAddr, "iacmaster/system")
	return &IaCRunner{
		Name:               name,
		Kind:               kind,
		MandatoryCommands:  mandatory_commands,
		Status:             models.InitStatus,
		artifactController: CreateIaCArtifactController(workingDir),
		Engine:             engine,
		CommandTimeout:     30,
		SystemPID:          systemPID,
		Error:              nil,
	}
}

func (l *IaCRunner) DeleteDeployment(deployment *msg.Deployment) bool {
	localPath := l.artifactController.TmpFolderPath + "/" + deployment.EnvironmentID + "/" + deployment.HomeFolder
	var err error
	if deployment.IaCArtifactType == "terraform" {
		err = l.terraformDestroy(localPath)
	} else if deployment.IaCArtifactType == "terragrunt" {
		err = l.terragruntDestroy(localPath)
	} else {
		err = errors.New("iac type not supported : " + deployment.IaCArtifactType)
	}

	if err != nil {
		deployment.Error = err.Error()
		return false
	}
	os.RemoveAll(localPath)
	return true
}
func (l *IaCRunner) SetDeployment(deployment *msg.Deployment) bool {
	l.Deployment = deployment
	if deployment.TerraformVersion != "" {
		err := l.setTerraformVersion(deployment.TerraformVersion)
		if err != nil {
			log.Println(err.Error())
			deployment.Error = err.Error()
			l.Status = models.FailedStatus
			return false
		}
	}
	l.Deployment = deployment
	err := l.GetRepo(deployment)
	if err != nil {
		log.Println("Error -> ", err.Error())
		deployment.Error = err.Error()
		l.Status = models.FailedStatus
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
	if deployment.IaCArtifactType == "terraform" {
		err = l.deployTerraformEnvironment(deployment)
	} else if deployment.IaCArtifactType == "terragrunt" {
		err = l.deployTerragruntEnvironment(deployment)
	} else {
		err = errors.New("iac artifact type not supported : " + deployment.IaCArtifactType)
	}

	if err != nil {
		deployment.Error = err.Error()
		log.Println("Deployment ", deployment.EnvironmentID, " deployment failed : ", err.Error())
		return false
	} else {
		log.Println("Deployment ", deployment.EnvironmentID, " deployment succeeded")
	}
	l.Status = models.SucceededStatus
	return true
}

func (l *IaCRunner) deployTerraformEnvironment(deployment *msg.Deployment) error {
	localPath := l.artifactController.TmpFolderPath + "/" + deployment.EnvironmentID + "/" + deployment.HomeFolder
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

func (l *IaCRunner) deployTerragruntEnvironment(deployment *msg.Deployment) error {
	localPath := l.artifactController.TmpFolderPath + "/" + deployment.EnvironmentID + "/" + deployment.HomeFolder
	err := l.terragruntInit(localPath)
	if err != nil {
		return err
	}
	err = l.terragruntApply(localPath, "", true)
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
			l.Status = models.FailedStatus
			l.Error = err
			return false
		}
	}
	return true
}

func (l *IaCRunner) SendLog(content string) {
	l.Engine.Send(l.SystemPID, &msg.Logging{Origin: l.Name, Content: content})
}

func (l *IaCRunner) runCommand(prog, running_dir string, commands []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(l.CommandTimeout))
	defer cancel()

	cmd := exec.CommandContext(ctx, prog, commands...)
	if running_dir != "" {
		cmd.Dir = running_dir
	}
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
		content := in.Text()
		log.Println(content)
		l.SendLog(content)
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("command timeout")
	}
	return nil
}

func (l *IaCRunner) setTerraformVersion(terraform_version string) error {
	prog := "tfenv"
	var commands []string
	commands = append(commands, "use")
	commands = append(commands, terraform_version)
	err := l.runCommand(prog, "", commands)
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

	err := l.runCommand(prog, "", commands)
	if err != nil {
		log.Println(err.Error())
	}

	commands = commands[:0]
	commands = append(commands, "account")
	commands = append(commands, "set")
	commands = append(commands, "--subscription")
	commands = append(commands, deployment.EnvironmentParameters["ARM_SUBSCRIPTION_ID"])

	err = l.runCommand(prog, "", commands)
	return err
}

func (l *IaCRunner) terraformInit(folder string) error {
	prog := "terraform"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "init")
	err := l.runCommand(prog, folder, commands)
	return err
}

func (l *IaCRunner) terragruntInit(folder string) error {
	prog := "terragrunt"
	var commands []string
	if folder != "" {
		commands = append(commands, "-chdir="+folder)
	}
	commands = append(commands, "--non-interactive")
	commands = append(commands, "run-all")
	commands = append(commands, "init")
	err := l.runCommand(prog, folder, commands)
	return err
}

func (l *IaCRunner) terraformPlan(folder, var_file_path string, save bool) error {
	prog := "terraform"
	var commands []string
	commands = append(commands, "plan")
	if save {
		commands = append(commands, "-out=plan.tfplan")
	}
	if var_file_path != "" {
		commands = append(commands, "-var-file="+var_file_path)
	}
	err := l.runCommand(prog, folder, commands)
	return err
}

func (l *IaCRunner) terraformApply(folder, var_file_path string, saved bool) error {
	prog := "terraform"
	var commands []string
	commands = append(commands, "apply")
	commands = append(commands, "-auto-approve")

	if saved {
		commands = append(commands, "plan.tfplan")
	}
	err := l.runCommand(prog, folder, commands)
	return err
}

func (l *IaCRunner) terragruntApply(folder, var_file_path string, saved bool) error {
	prog := "terragrunt"
	var commands []string
	commands = append(commands, "--non-interactive")
	commands = append(commands, "run-all")
	commands = append(commands, "apply")
	commands = append(commands, "-auto-approve")

	if saved {
		commands = append(commands, "plan.tfplan")
	}
	err := l.runCommand(prog, folder, commands)
	return err
}

func (l *IaCRunner) terraformDestroy(folder string) error {
	prog := "terraform"
	var commands []string
	commands = append(commands, "destroy")
	commands = append(commands, "-auto-approve")
	err := l.runCommand(prog, folder, commands)
	return err
}

func (l *IaCRunner) terragruntDestroy(folder string) error {
	prog := "terragrunt"
	var commands []string
	commands = append(commands, "--non-interactive")
	commands = append(commands, "run-all")
	commands = append(commands, "destroy")
	commands = append(commands, "-auto-approve")
	err := l.runCommand(prog, folder, commands)
	return err
}

func (s *IaCRunner) GetRunnerAddr() string {
	return fmt.Sprintf("%v:%v", os.Getenv("EXECUTOR_HOST_IP"), os.Getenv("RUNNER_HOST_PORT"))
}

func (s *IaCRunner) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		log.Println("Runner actor started on address -> ", ctx.Engine().Address())
		ctx.Send(s.SystemPID, &msg.RunnerStatus{Name: s.Name, Status: msg.Status_READY, Error: "", Operation: msg.OperationType_NO_OPERATION, Address: s.GetRunnerAddr()})
	case actor.Initialized:
		log.Println("Runner actor initialized")
	case *actor.PID:
		log.Println("Runner actor has god an ID")
	case *msg.Deployment:
		log.Println("Deployment object received")
		if m.Action == "create_env" {
			if !s.SetDeployment(m) {
				s.SendLog(s.Deployment.GetError())
				ctx.Send(s.SystemPID, &msg.RunnerStatus{Name: s.Name, Status: msg.Status_FAILED, Error: s.Deployment.GetError(), Operation: msg.OperationType_DEPLOYMENT, Address: s.GetRunnerAddr()})
				return
			}
			ctx.Send(s.SystemPID, &msg.RunnerStatus{Name: s.Name, Status: msg.Status_COMPLETED, Error: "", Operation: msg.OperationType_DEPLOYMENT, Address: s.GetRunnerAddr()})
		} else if m.Action == "destroy_env" {
			if !s.DeleteDeployment(m) {
				s.SendLog(s.Deployment.GetError())
				ctx.Send(s.SystemPID, &msg.RunnerStatus{Name: s.Name, Status: msg.Status_FAILED, Error: s.Deployment.GetError(), Operation: msg.OperationType_UNDEPLOYMENT, Address: s.GetRunnerAddr()})
				return
			}
			ctx.Send(s.SystemPID, &msg.RunnerStatus{Name: s.Name, Status: msg.Status_COMPLETED, Error: "", Operation: msg.OperationType_UNDEPLOYMENT, Address: s.GetRunnerAddr()})
		} else {
		}

	default:
		slog.Warn("server got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}
