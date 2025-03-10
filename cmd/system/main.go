package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/jdtotow/iacmaster/pkg/actors"
	"github.com/jdtotow/iacmaster/pkg/api"
	"github.com/jdtotow/iacmaster/pkg/initializers"
	"github.com/jdtotow/iacmaster/pkg/models"
)

func init() {
	initializers.LoadEnvVariables()
	//creating tmp folder
	pwd, _ := os.Getwd()
	path := pwd + "/tmp"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir|0755)
	}
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print("[IaCMaster] " + time.Now().UTC().Format("[02-01-2006 - 15:04:05] ") + string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	channel := make(chan models.HTTPMessage)

	http_server := api.CreateSystemServer(&channel)
	go http_server.Start()
	system_port := os.Getenv("IACMASTER_SYSTEM_PORT")
	system_address := os.Getenv("IACMASTER_SYSTEM_ADDRESS")
	if system_address == "" {
		system_address = "0.0.0.0"
	}
	if system_port == "" {
		system_port = "3434"
	}
	r := remote.New(system_address+":"+system_port, remote.NewConfig())
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		log.Fatal("failed to create engine for iacmaster system", "error", err)
	}
	pid := engine.Spawn(actors.CreateSystemActor(&channel), "iacmaster", actor.WithID("system"))
	log.Println("System pid -> ", pid)
	for range 10 {
		log.Println("waiting for msg")
		time.Sleep(time.Second * 10)
	}
}
