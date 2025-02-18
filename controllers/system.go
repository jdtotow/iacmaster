package controllers

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/jdtotow/iacmaster/models"
)

type System struct {
	node         *models.Node
	dbController *DBController
	seController *SecurityController
	channel      *chan models.HTTPMessage
	peers        []*models.Node
}

func CreatePeers(settings, myName string) []*models.Node {
	//settings="node_name1=addr1,node_name2=addr2"
	result := []*models.Node{}
	if len(settings) == 0 {
		return result
	}
	for _, chunk := range strings.Split(settings, ",") {
		setting := strings.Split(chunk, "=")
		if setting[0] == myName {
			continue //skipping myself
		}
		n := &models.Node{
			Type:   models.NodeType("primary"),
			Name:   setting[0],
			Addr:   setting[1],
			Status: models.NodeStatus("unknown"),
		}
		result = append(result, n)
	}
	log.Printf("%v nodes found in settings\n", len(result))
	return result
}
func CreateSystem(channel *chan models.HTTPMessage) *System {
	var nodeName string = os.Getenv("NODE_NAME")
	var nodeType string = os.Getenv("NODE_TYPE")
	var clusterSetting string = os.Getenv("CLUSTER")

	if nodeName == "" {
		log.Fatal("Please set NODE_NAME variable")
	}
	if nodeType == "" {
		log.Fatal("Please set NODE_TYPE variable")
	}

	n := &models.Node{
		Type:   models.NodeType(nodeType),
		Name:   nodeName,
		Mode:   models.NodeMode("standalone"),
		Status: models.NodeStatus("init"),
	}
	return &System{
		node:         n,
		dbController: CreateDBController(),
		seController: CreateSecurityController(),
		channel:      channel,
		peers:        CreatePeers(clusterSetting, nodeName),
	}
}
func (s *System) UpdateTableSchema() {
	s.dbController.db_client.AutoMigrate(
		&models.User{},
		&models.Role{},
		//&models.Project{},
		&models.Organization{},
		//&models.IaCExecutionSettings{},
		//&models.IaCArtifact{},
		&models.Group{},
		//&models.Environment{},
		//&models.CloudCredential{},
	)
}
func (s *System) CheckMandatoryTableAndData() error {
	//Verify the system organization
	org := models.Organization{
		Name: "system",
	}
	result := s.dbController.db_client.Create(&org)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	systemUser := models.User{
		Email:          "system@iacmaster",
		Fullname:       "IaCMaster",
		Username:       "iacmaster",
		Password:       "--",
		OrganizationID: org.ID,
	}
	result = s.dbController.CreateInstance(&systemUser)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	//verify the system role
	role := models.Role{
		Name:   "system",
		UserID: systemUser.ID,
	}
	result = s.dbController.CreateInstance(&role)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	} else {
		log.Println("Role system stored")
	}
	group := models.Group{
		Name: "system",
	}
	result = s.dbController.CreateInstance(&group)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	s.dbController.db_client.Model(&systemUser).Association("Groups").Append(&group)

	//
	return nil
}

func (s *System) Start() {
	if s.node.Type == models.Primary {
		//s.UpdateTableSchema()
		//err := s.CheckMandatoryTableAndData()
		//if err != nil {
		//	log.Fatal("Cannot continue, missing mandatory data")
		//}
		//
		log.Println("IaC Master logic started !")
		var message models.HTTPMessage
		for {
			log.Println("Waiting for event ...")
			message = <-*s.channel
			s.Handle(message)
			time.Sleep(time.Second)
		}
	}
}

func (s *System) Handle(message models.HTTPMessage) {
	log.Println("message -> ", message)
}
