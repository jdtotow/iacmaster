package main

import (
	"fmt"
	"github.com/jdtotow/iacmaster/entities" entities 
)

func main() {
	fmt.Println("welcome to IaC Master")
	org := entities.CreateOrganization("Swisscom")
}
