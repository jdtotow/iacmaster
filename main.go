package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jdtotow/iacmaster/api"
	"github.com/jdtotow/iacmaster/controllers"
	"github.com/jdtotow/iacmaster/initializers"
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

	fmt.Println("Welcome to IaC Master\nStarting api server ...")
	dbController := controllers.CreateDBController(dbUri)
	seController := controllers.CreateSecurityController(secretKey, dbController)
	system := controllers.CreateSystem(nodeType, nodeName, dbController)
	system.Start()
	http_server := api.CreateServer(port, dbController, seController)
	http_server.Start()
}
