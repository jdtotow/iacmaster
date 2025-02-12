package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jdtotow/iacmaster/api"
	"github.com/jdtotow/iacmaster/controllers"
	"github.com/jdtotow/iacmaster/initializers"
	"github.com/jdtotow/iacmaster/models"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	var port int
	port, _ = strconv.Atoi(os.Getenv("API_PORT"))
	var dbUri string = os.Getenv("DB_URI")
	var secretKey string = os.Getenv("SECRET_KEY")
	var nodeName string = os.Getenv("NODE_NAME")
	var nodeType string = os.Getenv("NODE_TYPE")

	channel := make(chan models.HTTPMessage)

	fmt.Println("Initializing controllers ...")
	dbController := controllers.CreateDBController(dbUri)
	seController := controllers.CreateSecurityController(secretKey)
	http_server := api.CreateServer(port, &channel)
	system := controllers.CreateSystem(nodeType, nodeName, dbController, seController, http_server, &channel)

	system.Start()
}
