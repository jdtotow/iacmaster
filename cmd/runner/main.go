package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/jdtotow/iacmaster/pkg/actors"
	"github.com/jdtotow/iacmaster/pkg/models"
)

func main() {
	executor_name := os.Getenv("EXECUTOR_NAME")
	mandatory_commands_str := os.Getenv("MANDATORY_COMMANDS")
	working_dir := os.Getenv("WORKING_DIR")
	kind := os.Getenv("EXECUTOR_KIND")

	mandatory_commands := []string{}
	mandatory_commands = append(mandatory_commands, strings.Split(mandatory_commands_str, ",")...)

	ifaces, err := net.Interfaces()
	// handle err
	var private_ip string
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsPrivate() {
				private_ip = ip.String()
			}
			// process IP address
		}
	}

	r := remote.New(private_ip+":8787", remote.NewConfig())
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		log.Fatal("failed to create engine for runner", "error", err)
	}
	engine.Spawn(actors.CreateRunnerActor(working_dir, executor_name, mandatory_commands, models.ExecutorKind(kind)), "iacmaster", actor.WithID("runner"))
}
