package main

import (
	"github.com/jdtotow/iacmaster/pkg/api"
	"github.com/jdtotow/iacmaster/pkg/initializers"
)

func init() {
	initializers.LoadEnvVariables()
}
func main() {
	server := api.CreateServiceServer()
	server.Start()
}
