package controllers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
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
func (s *System) CheckMandatoryTableAndData() error {
	//Verify the system organization
	org := models.Organization{}
	s.dbController.db_client.AutoMigrate(&org)
	s.dbController.db_client.First(&org, "id = ? ", 1)
	if org.GetName() != "system" {
		fmt.Println("The system organization does not exist, it will be created")
		org.SetName("system")
		org.SetUuid(uuid.NewString())
		result := s.dbController.db_client.Create(&org)
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
	}
	//verify the system role
	role := models.Role{}
	s.dbController.db_client.AutoMigrate(&role)
	s.dbController.db_client.First(&role, "id = ?", 1)
	if role.GetName() != "system" {
		log.Println("The system role does not exist, it will be created")
		role.SetName("system")
		role.SetUuid(uuid.NewString())
		result := s.dbController.CreateInstance(&role)
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
		//create all roles
		roles := []string{
			"owner",
			"administrator",
			"project_manager",
			"deployer",
			"reader",
		}
		for _, _role := range roles {
			r := models.Role{
				Name: _role,
				Uuid: uuid.NewString(),
			}
			s.dbController.CreateInstance(&r)
		}
	}
	//verify the usergroup
	group := models.UserGroup{}
	s.dbController.db_client.AutoMigrate(&group)
	s.dbController.db_client.First(&group, "id = ?", 1)
	if group.GetName() != "system" {
		log.Println("The system role does not exist, it will be created")
		group.SetName("system")
		group.SetUuid(uuid.NewString())
		result := s.dbController.CreateInstance(&group)
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
	}
	//verify the system user
	systemUser := models.User{}
	s.dbController.db_client.AutoMigrate(&models.User{})
	s.dbController.db_client.First(&systemUser, "id = ?", 1)
	if systemUser.GetUsername() != "iacmaster" {
		log.Println("The system user does not exist, it will be created")
		systemUser.SetEmail("iacmaster@system")
		systemUser.SetFullname("IaCMaster")
		systemUser.SetPassword("")
		systemUser.SetUsername("iacmaster")
		systemUser.AddRole(role)
		systemUser.SetOrganization(org)
		systemUser.SetUuid(uuid.NewString())
		result := s.dbController.CreateInstance(&systemUser)
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
	}
	//
	return nil
}

func (s *System) Start() {
	if s.node.Type == models.Primary {
		err := s.CheckMandatoryTableAndData()
		if err != nil {
			log.Fatal("Cannot continue, missing mandatory data")
		}
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
