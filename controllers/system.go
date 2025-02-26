package controllers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jdtotow/iacmaster/models"
	"github.com/jdtotow/iacmaster/worker"
)

type System struct {
	node               *models.Node
	dbController       *DBController
	seController       *SecurityController
	artifactController *IaCArtifactController
	channel            *chan models.HTTPMessage
	peers              []*models.Node
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
		node:               n,
		dbController:       CreateDBController(),
		seController:       CreateSecurityController(),
		artifactController: CreateIaCArtifactController("./tmp"),
		channel:            channel,
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
		log.Println("Preparing deployment with ", env.IaCExecutionSettings.TerraformVersion)
		err := s.artifactController.GetRepo(env.IaCArtifact.ScmUrl, env.IaCExecutionSettings.Token.Token, env.IaCExecutionSettings.Token.Username, env.IaCArtifact.Revision, env.IaCArtifact.ProxyUrl, env.IaCArtifact.ProxyUsername, env.IaCArtifact.ProxyPassword, message.Metadata["object_id"])
		if err != nil {
			log.Println("Could not clone git repo")
			return
		}
		//moving variables.tfvars
		pwd, _ := os.Getwd()
		sourcePath := pwd + "/tmp/" + message.Metadata["object_id"] + ".tfvars"
		if _, err := os.Stat(sourcePath); err == nil {
			// File exists, proceed to move
			err = os.Rename(sourcePath, pwd+"/tmp/"+message.Metadata["object_id"]+"/"+env.IaCArtifact.HomeFolder+"/variables.tfvars")
			if err != nil {
				log.Printf("Error moving file: %v\n", err)
				return
			}
			log.Println("File moved successfully.")
		}
		// create worker
		runner := &Runner{}
		docker_worker := runner.Create("default", "docker")
		docker_image := ""
		if env.IaCArtifact.Type == "terraform" {
			docker_image = "iacmaster_worker:latest"
		}
		info := worker.JobData{
			VolumePath:            pwd + "/tmp/" + message.Metadata["object_id"] + "/" + env.IaCArtifact.HomeFolder,
			EnvironmentID:         message.Metadata["object_id"],
			EnvironmentParameters: cloud_credential.Variables,
			DockerImage:           docker_image,
			TerraformVersion:      env.IaCExecutionSettings.TerraformVersion,
			WorkingDir:            "/tmp/" + message.Metadata["object_id"] + "/" + env.IaCArtifact.HomeFolder,
		}
		docker_worker.SetJobInfo(info)
	}
}
