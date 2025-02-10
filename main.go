package main

import (
	"fmt"

	"github.com/jdtotow/iacmaster/api"
)

func main() {
	var port int = 3000
	fmt.Println("Welcome to IaC Master\nStartinh api server ...")
	http_server := api.CreateServer(port)
	http_server.Start()
}
