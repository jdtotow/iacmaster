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

	fmt.Println("Welcome to IaC Master\nStartinh api server ...")
	dbController := controllers.CreateDBController(dbUri)
	seController := controllers.CreateSecurityController(secretKey)
	http_server := api.CreateServer(port, dbController, seController)
	http_server.Start()
}
