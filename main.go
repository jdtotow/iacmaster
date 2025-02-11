package main

import (
	"fmt"
	"strconv"

	"github.com/jdtotow/iacmaster/api"
	"github.com/joho/godotenv"
)

func main() {
	var port int = 3000
	var dbPort int = 27017
	var dbUsername string = ""
	var dbPassword string = ""
	var dbUri string = "mongodb://localhost:27017"
	var dbName string = "iacmaster"
	var secretKey string = "your-secret-key"
	//Reading variables from .env file
	envFile, err := godotenv.Read(".env")
	if err == nil {
		port, _ = strconv.Atoi(envFile["API_PORT"])
		dbPort, _ = strconv.Atoi(envFile["DB_PORT"])
		dbUsername = envFile["DB_USERNAME"]
		dbPassword = envFile["DB_PASSWORD"]
		dbUri = envFile["DB_URI"]
		dbName = envFile["DB_NAME"]
		secretKey = envFile["SECRET_KEY "]
	}

	fmt.Println("Welcome to IaC Master\nStartinh api server ...")
	http_server := api.CreateServer(dbUri, dbUsername, dbPassword, dbName, secretKey, port, dbPort)
	http_server.Start()
}
