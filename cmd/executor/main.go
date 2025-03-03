package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/jdtotow/iacmaster/pkg/controllers"
	"github.com/jdtotow/iacmaster/pkg/models"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Deployment file not provided")
	}
	deploymentFile := os.Args[1]
	file, err := os.Open(deploymentFile)
	if err != nil {
		log.Fatalf("Failed to open task file: %v", err)
	}
	defer file.Close()
	executor_name := os.Getenv("EXECUTOR_NAME")
	mandatory_commands_str := os.Getenv("MANDATORY_COMMANDS")
	working_dir := os.Getenv("WORKING_DIR")
	kind := os.Getenv("EXECUTOR_KIND")

	mandatory_commands := []string{}
	mandatory_commands = append(mandatory_commands, strings.Split(mandatory_commands_str, ",")...)
	executor := controllers.CreateIaCExecutor(working_dir, executor_name, mandatory_commands, controllers.ExecutorKind(kind))

	var deployment models.Deployment
	if err := json.NewDecoder(file).Decode(&deployment); err != nil {
		log.Fatalf("Failed to parse task file: %v", err)
		executor.State.Status = controllers.FailedStatus
		executor.State.Error = err
	}
	executor.SetDeployment(&deployment)
}
