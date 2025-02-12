package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jdtotow/iacmaster/models"
)

type System struct {
	node         *models.Node
	dbController *DBController
	seController *SecurityController
}

func CreateSystem(nodeType, nodeName string, dbController *DBController, seController *SecurityController) *System {
	n := &models.Node{
		Type:   models.NodeType(nodeType),
		Name:   nodeName,
		Status: models.NodeStatus("init"),
	}
	return &System{
		node:         n,
		dbController: dbController,
		seController: seController,
	}
}
func (s *System) CheckMandatoryTableAndData() error {
	org := models.Organization{}
	err := s.dbController.db_client.AutoMigrate(&models.Organization{})
	fmt.Println("Model creation", err)
	s.dbController.db_client.First(&org, "id = ? ", 1)
	if org.GetName() != "system" {
		fmt.Println("The system organization does not exist, it will be created")
		org.SetName("system")
		result := s.dbController.db_client.Create(&org)
		if result.Error != nil {
			fmt.Println(result.Error)
			return result.Error
		}
	}
	systemUser := models.User{}
	s.dbController.db_client.AutoMigrate(&models.User{})
	s.dbController.db_client.First(&systemUser, "id = ?", 1)
	if systemUser.GetUsername() != "iacmaster" {
		fmt.Println("The system user does not exist, it will be created")
		systemUser.SetEmail("iacmaster@system")
		systemUser.SetFullname("IaCMaster")
		systemUser.SetPassword("")
		systemUser.SetUsername("iacmaster")
		result := s.dbController.db_client.Create(&systemUser)
		if result.Error != nil {
			fmt.Println(result.Error)
			return result.Error
		}
	}
	return nil
}

func (s *System) Start() {
	if s.node.Type == models.Primary {
		err := s.CheckMandatoryTableAndData()
		if err != nil {
			log.Fatal("Cannot continue, missing mandatory data")
		}
	}
}

func (s *System) Handle(context *gin.Context, objectName string) {
	fmt.Println(objectName, " called ")
}
