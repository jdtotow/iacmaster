package main

import (
	"log"
	"os"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/jdtotow/iacmaster/pkg/actors"
	"github.com/jdtotow/iacmaster/pkg/models"
)

func main() {
	executor_name := os.Getenv("DEPLOYMENT_ID")
	mandatory_commands_str := os.Getenv("MANDATORY_COMMANDS")
	working_dir := os.Getenv("RUNNER_WORKING_DIR")
	kind := os.Getenv("EXECUTOR_KIND")

	mandatory_commands := []string{}
	mandatory_commands = append(mandatory_commands, strings.Split(mandatory_commands_str, ",")...)
	private_ip := os.Getenv("EXECUTOR_HOST_IP")
	runner_host_port := os.Getenv("RUNNER_HOST_PORT")
	r := remote.New(private_ip+":"+runner_host_port, remote.NewConfig())
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		log.Fatal("failed to create engine for runner", "error", err)
	}
	engine.Spawn(actors.CreateRunnerActor(working_dir, executor_name, mandatory_commands, models.ExecutorKind(kind), engine), "runner", actor.WithID(os.Getenv("DEPLOYMENT_ID")))
	select {}
}
