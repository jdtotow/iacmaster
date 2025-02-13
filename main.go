package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	var port int
	port, _ = strconv.Atoi(os.Getenv("API_PORT"))
	var dbUri string = os.Getenv("DB_URI")
	var secretKey string = os.Getenv("SECRET_KEY")
	var nodeName string = os.Getenv("NODE_NAME")
	var nodeType string = os.Getenv("NODE_TYPE")
	var clusterSetting string = os.Getenv("CLUSTER")

	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	channel := make(chan models.HTTPMessage)

	http_server := api.CreateServer(port, &channel)
	system := controllers.CreateSystem(nodeType, nodeName, http_server, &channel, dbUri, secretKey, clusterSetting)

	go http_server.Start()
	system.Start()
}
