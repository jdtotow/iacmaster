package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"slices"

	"github.com/jdtotow/iacmaster/pkg/models"
)

type System struct {
	node               *models.Node
	dbController       *DBController
	seController       *SecurityController
	artifactController *IaCArtifactController
	channel            *chan models.HTTPMessage
	peers              []*models.Node
	serviceUrl         string
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
	nodeMode := "standalone"
	var nodeAttributes []models.NodeAttribute
	if nodeMode == "standalone" {
		nodeAttributes = []models.NodeAttribute{
			models.NodeAttribute("log_event"),
			models.NodeAttribute("manager"),
			models.NodeAttribute("executor"),
		}

	}

	n := &models.Node{
		Type:       models.NodeType(nodeType),
		Name:       nodeName,
		Mode:       models.NodeMode("standalone"),
		Status:     models.NodeStatus("init"),
		Attributes: nodeAttributes,
	}
	return &System{
		node:               n,
		dbController:       CreateDBController(),
		seController:       CreateSecurityController(),
		artifactController: CreateIaCArtifactController("./tmp"),
		channel:            channel,
		serviceUrl:         os.Getenv("SERVICE_URL"),
		peers:              CreatePeers(clusterSetting, nodeName),
	}
}
func (s *System) UpdateTableSchema() {
	s.dbController.db_client.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Project{},
		&models.Token{},
		&models.IaCSystem{},
		&models.Organization{},
		&models.IaCExecutionSettings{},
		&models.IaCArtifact{},
		&models.Group{},
		&models.Environment{},
		&models.CloudCredential{},
	)
}
func (s *System) CreateTablesAndMandatoryData() error {
	if s.CheckMandatoryTableAndData() {
		return nil
	}
	s.UpdateTableSchema()
	sys := models.IaCSystem{
		Name: "IaCSystem",
	}
	result := s.dbController.db_client.Create(&sys)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	org := models.Organization{
		Name: "system",
	}
	result = s.dbController.db_client.Create(&org)
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
func (s *System) CheckMandatoryTableAndData() bool {
	return s.dbController.db_client.Migrator().HasTable("ia_c_systems")
}

func (s *System) IsNodeManager() bool {
	return slices.Contains(s.node.Attributes, models.NodeAttribute("manager"))
}
func (s *System) IsNodeEventLog() bool {
	return slices.Contains(s.node.Attributes, models.NodeAttribute("log_event"))
}
func (s *System) IsNodeExecutor() bool {
	return slices.Contains(s.node.Attributes, models.NodeAttribute("executor"))
}

func (s *System) Start() {
	if s.node.Type == models.Primary {

		err := s.CreateTablesAndMandatoryData()
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
	if message.Metadata["action"] == "create_env" {
		env := models.Environment{}
		s.dbController.GetClient().Preload("Project").Preload("IaCArtifact").Preload("IaCExecutionSettings").First(&env, "id = ?", message.Metadata["object_id"])
		cloud_credential := models.CloudCredential{}
		env_settings_id := fmt.Sprintf("%v", env.IaCExecutionSettings.CloudCredentialID)
		result := s.dbController.GetObjectByID(&cloud_credential, env_settings_id)
		if result.Error != nil {
			log.Fatal("Could not retrieve coud credential")
		}
		if env.Name == "" {
			log.Println("Object not found")
			return
		}
		// send request to service
		deployment := models.Deployment{}
		deployment.CloudDestination = string(env.IaCExecutionSettings.DestinationCloud)
		all_env_parameters := map[string]string{}
		for k, v := range cloud_credential.Variables {
			all_env_parameters[k] = v
		}
		for k, v := range env.Project.Variables {
			all_env_parameters[k] = v
		}
		deployment.EnvironmentParameters = all_env_parameters
		deployment.EnvironmentID = message.Metadata["object_id"]
		deployment.HomeFolder = env.IaCArtifact.HomeFolder
		deployment.GitData.Url = env.IaCArtifact.ScmUrl
		deployment.GitData.Revision = env.IaCArtifact.Revision
		deployment.GitData.ProxyUrl = env.IaCArtifact.ProxyUrl
		deployment.GitData.ProxyUsername = env.IaCArtifact.ProxyUsername
		deployment.GitData.ProxyPassword = env.IaCExecutionSettings.Token.Token
		deployment.TerraformVersion = env.IaCExecutionSettings.TerraformVersion

		deploy_json, err := json.Marshal(deployment)
		if err != nil {
			fmt.Println("Could not send serialize deployment object : ", err.Error())
		} else {
			resp, err := http.Post(s.serviceUrl+"/deployment", "application/json", bytes.NewBuffer(deploy_json))
			log.Println("Request sent to service")
			if err == nil {
				fmt.Println(resp.StatusCode)
			} else {
				fmt.Println(err.Error())
			}
		}
	} else if message.Metadata["action"] == "destroy_env" {
		// send request to service
		deployment := models.Deployment{}
		env := models.Environment{}
		s.dbController.GetClient().Preload("Project").Preload("IaCArtifact").Preload("IaCExecutionSettings").First(&env, "id = ?", message.Metadata["object_id"])
		deployment.EnvironmentID = message.Metadata["object_id"]
		deployment.HomeFolder = env.IaCArtifact.HomeFolder
		deploy_json, err := json.Marshal(deployment)
		if err != nil {
			fmt.Println("Could not send serialized deployment object : ", err.Error())
		} else {
			resp, err := http.Post(s.serviceUrl+"/destroy", "application/json", bytes.NewBuffer(deploy_json))
			log.Println("Request sent to service")
			if err == nil {
				fmt.Println(resp.StatusCode)
			} else {
				fmt.Println(err.Error())
			}
		}
	} else {
		fmt.Println("Unknown action: ", message.Metadata["action"])
	}
}
