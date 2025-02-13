package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jdtotow/iacmaster/api"
	"github.com/jdtotow/iacmaster/controllers"
	"github.com/jdtotow/iacmaster/initializers"
	"github.com/jdtotow/iacmaster/models"
)

func init() {
	initializers.LoadEnvVariables()
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

	http_server := api.CreateServer(&channel)
	system := controllers.CreateSystem(&channel)

	go http_server.Start()
	system.Start()
}
