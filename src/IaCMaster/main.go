package main

import (
	"fmt"

	"github.com/jdtotow/iacmaster/entities"
)

func main() {
	fmt.Println("welcome to IaC Master")
	org := entities.CreateOrganization("Swisscom")
	fmt.Println(org)
}
