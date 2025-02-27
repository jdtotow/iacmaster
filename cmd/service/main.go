package main

import (
	"workerservice/api"
	"workerservice/initializers"
)

func init() {
	initializers.LoadEnvVariables()
}
func main() {
	server := api.CreateServer()
	server.Start()
}
