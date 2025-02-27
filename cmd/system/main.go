package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jdtotow/iacmaster/pkg/api"
	"github.com/jdtotow/iacmaster/pkg/controllers"
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
	system := controllers.CreateSystem(&channel)

	go http_server.Start()
	system.Start()
}
