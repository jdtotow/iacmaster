package main

import (
	"fmt"

	"github.com/jdtotow/iacmaster/src/IaCMaster/entities"
)

func main() {
	fmt.Println("welcome to IaC Master")
	org := entities.CreateOrganization("Swisscom")
	fmt.Println(org)
}
