package main

import (
	"fmt"

	"github.com/jdtotow/iacmaster/api"
)

func main() {
	var port int = 3000
	var dbPort int = 27017
	var dbUsername string = ""
	var dbPassword string = ""
	var dbUri string = "mongodb://localhost"
	var dbName string = "iacmaster"

	fmt.Println("Welcome to IaC Master\nStartinh api server ...")
	http_server := api.CreateServer(dbUri, dbUsername, dbPassword, dbName, port, dbPort)
	http_server.Start()
}
