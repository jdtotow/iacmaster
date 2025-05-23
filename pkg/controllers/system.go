package controllers

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"maps"

	"slices"

	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type System struct {
	dbController       *DBController
	seController       *SecurityController
	artifactController *IaCArtifactController
	serviceUrl         string
	executorManager    *ExecutorManager
	nodeInfo           *msg.NodeInfo
	attributes         []models.NodeAttribute
	started            bool
	ActorEngine        *actor.Engine
	EventHubActorPID   *actor.PID
	ClusterInfo        *models.ClusterInfo
}

func CreateSystem() *System {
	var executionPlatform string = os.Getenv("EXECUTION_PLATFORM")
	pwd, _ := os.Getwd()
	working_dir := pwd + "/tmp"
	return &System{
		dbController:       CreateDBController(),
		seController:       CreateSecurityController(),
		artifactController: CreateIaCArtifactController("./tmp"),
		serviceUrl:         os.Getenv("SERVICE_URL"),
		executorManager:    CreateExecutorManager(working_dir, executionPlatform),
		attributes:         []models.NodeAttribute{},
		started:            false,
	}
}

func CreateEventHubActor() actor.Producer {
	return func() actor.Receiver {
		return CreateEventHub()
	}
}

func (s *System) UpdateTableSchema() {
	s.dbController.Db_client.AutoMigrate(
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
		&models.Access{},
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
	result := s.dbController.Db_client.Create(&sys)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	org := models.Organization{
		Name: "system",
	}
	result = s.dbController.Db_client.Create(&org)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	systemUser := models.User{
		Email:    "system@iacmaster",
		Fullname: "IaCMaster",
		Username: "iacmaster",
		Password: "--",
	}
	result = s.dbController.CreateInstance(&systemUser)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	s.dbController.Db_client.Model(&systemUser).Association("Organizations").Append(&org)

	//verify the system role
	role := models.Role{
		Name: "system",
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
	s.dbController.Db_client.Model(&systemUser).Association("Groups").Append(&group)
	return nil
}
func (s *System) CheckMandatoryTableAndData() bool {
	return s.dbController.Db_client.Migrator().HasTable("ia_c_systems")
}

func (s *System) IsNodeManager() bool {
	return slices.Contains(s.attributes, models.ManagerNodeAttribute)
}
func (s *System) IsNodeEventLog() bool {
	return slices.Contains(s.attributes, models.LoggingNodeAttribute)
}
func (s *System) IsNodeExecutor() bool {
	return slices.Contains(s.attributes, models.ExecutorNodeAttribute)
}

func (s *System) Start() {
	if s.started {
		return
	}
	if os.Getenv("CLUSTER") == "" {
		s.attributes = append(s.attributes, models.ExecutorNodeAttribute)
		s.attributes = append(s.attributes, models.ManagerNodeAttribute)
		s.attributes = append(s.attributes, models.LoggingNodeAttribute)
		s.attributes = append(s.attributes, models.APINodeAttribute)
	} else {

		s.ClusterInfo = models.CreateClusterInfo("iacmaster")

		for _, chunk := range strings.Split(os.Getenv("CLUSTER"), ",") {
			setting := strings.Split(chunk, "=")
			nodeName := setting[0]
			if os.Getenv("NODE_NAME") == nodeName {
				continue //skipping myself
			}
			nodeAddr := setting[1] + ":3000"
			node := &msg.NodeInfo{
				Name:        nodeName,
				Addr:        nodeAddr,
				NodeType:    uint32(models.Secondary),
				Deployments: []string{},
			}
			s.ClusterInfo.AddNode(node)
		}
		log.Println("Cluster info created")
	}

	if s.nodeInfo.NodeType == uint32(models.Primary) {
		err := s.CreateTablesAndMandatoryData()
		if err != nil {
			log.Fatal("Cannot continue, missing mandatory data")
		}
		s.attributes = append(s.attributes, models.ManagerNodeAttribute)
		s.attributes = append(s.attributes, models.LoggingNodeAttribute)
		s.attributes = append(s.attributes, models.APINodeAttribute)
	} else {
		s.attributes = append(s.attributes, models.ExecutorNodeAttribute)
	}
	log.Println("IaC Master logic started !")
	s.started = true
	//sending node attribute message to api
	_msg := &msg.NodeAttribute{
		NodeName: s.nodeInfo.Name,
	}
	for _, attribute := range s.attributes {
		if attribute == models.ExecutorNodeAttribute {
			_msg.Attribute = append(_msg.Attribute, msg.NodeAttributeType_EXECUTOR)
		}
		if attribute == models.ManagerNodeAttribute {
			_msg.Attribute = append(_msg.Attribute, msg.NodeAttributeType_MANAGER)
		}
		if attribute == models.LoggingNodeAttribute {
			_msg.Attribute = append(_msg.Attribute, msg.NodeAttributeType_LOGGING)
		}
		if attribute == models.APINodeAttribute {
			_msg.Attribute = append(_msg.Attribute, msg.NodeAttributeType_API)
		}
	}
	senderAddr := os.Getenv("IACMASTER_SYSTEM_ADDRESS") + ":" + os.Getenv("IACMASTER_SYSTEM_PORT")
	apiPID := actor.NewPID(senderAddr, "iacmaster/api")
	s.ActorEngine.Send(apiPID, _msg)
}

func (s *System) SendSubscriptionRequest(subscriber_id string, subscriptions []*msg.Subscription) {
	req := &msg.SubscriptionRequest{
		Id:           subscriber_id,
		Subcriptions: subscriptions,
	}
	s.ActorEngine.Send(s.EventHubActorPID, req)
}
func (s *System) Subscribe() {
	subscriptions := []*msg.Subscription{
		{
			ActionType:   "actorengineaction",
			Destination:  s.nodeInfo.Addr + ":iacmaster/system",
			DeploymentID: "",
			EventType:    "log",
		},
	}
	s.SendSubscriptionRequest("system", subscriptions)
}

func (s *System) Handle(operation *msg.Operation) {
	log.Println("message -> ", operation)
	if operation.Action == "create_env" {
		env := models.Environment{}
		s.dbController.GetClient().Preload("Project").Preload("IaCArtifact").Preload("IaCExecutionSettings").First(&env, "id = ?", operation.ObjectID)
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
		deployment := msg.Deployment{}
		git_data := msg.GitData{}
		deployment.Action = operation.Action
		deployment.CloudDestination = string(env.IaCExecutionSettings.DestinationCloud)
		all_env_parameters := map[string]string{}
		maps.Copy(all_env_parameters, cloud_credential.Variables)
		maps.Copy(all_env_parameters, env.Project.Variables)

		deployment.EnvironmentParameters = all_env_parameters
		deployment.EnvironmentID = operation.ObjectID
		deployment.HomeFolder = env.IaCArtifact.HomeFolder
		deployment.IaCArtifactType = env.IaCArtifact.Type
		deployment.DetectDrift = env.IaCExecutionSettings.DetectDrift
		deployment.AutoRedeployOnGitChange = env.IaCExecutionSettings.AutoRedeployOnGitChange
		git_data.Url = env.IaCArtifact.ScmUrl
		git_data.Revision = env.IaCArtifact.Revision
		git_data.ProxyUrl = env.IaCArtifact.ProxyUrl
		git_data.ProxyUsername = env.IaCArtifact.ProxyUsername
		git_data.ProxyPassword = env.IaCExecutionSettings.Token.Token
		deployment.TerraformVersion = env.IaCExecutionSettings.TerraformVersion

		deployment.GitData = &git_data
		if s.IsNodeExecutor() {
			err := s.executorManager.StartDeployment(&deployment)
			if err != nil {
				log.Println(err)
			} else {
				s.ClusterInfo.AddDeploymentToNode(s.nodeInfo.Name, deployment.EnvironmentID)
			}
		} else {
			log.Println("Deployment request must be sent to node executor")
			random_node := s.ClusterInfo.GetRandomNode()
			destinationPID := s.ClusterInfo.GetPIDNode(random_node.Name, "system")
			s.ActorEngine.Send(destinationPID, &deployment)
			s.ClusterInfo.AddDeploymentToNode(random_node.Name, deployment.EnvironmentID)
		}

	} else if operation.Action == "destroy_env" {
		// send request to service
		deployment := &msg.Deployment{}
		env := models.Environment{}
		s.dbController.GetClient().Preload("Project").Preload("IaCArtifact").Preload("IaCExecutionSettings").First(&env, "id = ?", operation.ObjectID)
		deployment.Action = operation.Action
		deployment.EnvironmentID = operation.ObjectID
		deployment.HomeFolder = env.IaCArtifact.HomeFolder
		deployment.IaCArtifactType = env.IaCArtifact.Type
		// getting node executing
		node := s.ClusterInfo.GetNodeExecutingDeployment(deployment.EnvironmentID)
		if node == nil {
			err := s.executorManager.StartDeployment(deployment)
			if err != nil {
				log.Println(err)
			}
		} else {
			destinationPID := s.ClusterInfo.GetPIDNode(node.Name, "system")
			s.ActorEngine.Send(destinationPID, &deployment)
		}
	} else {
		fmt.Println("Unknown action: ", operation.ObjectID)
	}
}

func (s *System) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		log.Println("System actor started at -> ", ctx.Engine().Address())
		s.ActorEngine = ctx.Engine()
	case actor.Stopped:
		log.Println("System actor has stopped")
	case *msg.Operation:
		s.Handle(m)
	case actor.Initialized:
		log.Println("System actor initialized")
	case *actor.PID:
		log.Println("System actor has god an ID")
	case *msg.RunnerStatus:
		log.Println("Runner with id ", m.Name, " is ", m.Status)
		s.HandlerRunnerStatus(m, ctx)
	case *msg.Logging:
		log.Println("[", m.Origin, "] ", m.Content)
	case *msg.Deployment:
		log.Println("Deployment object received")
		err := s.executorManager.StartDeployment(m)
		if err != nil {
			log.Println(err)
		}
	case *msg.NodeInfo:
		s.nodeInfo = m
		log.Println("System has received node info message received -> ", m)
		if m.NodeStatus == uint32(models.Running) {
			s.Start()
		}
		if s.nodeInfo.NodeType == uint32(models.Primary) {
			s.EventHubActorPID = ctx.Engine().Spawn(CreateEventHubActor(), "iacmaster", actor.WithID("eventhub"))
			if s.EventHubActorPID != nil {
				s.Subscribe()
			}

		} else {
			if s.EventHubActorPID != nil {
				ctx.Engine().Stop(s.EventHubActorPID)
			}
		}

	default:
		slog.Warn("server got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}

func (s *System) HandlerRunnerStatus(runner *msg.RunnerStatus, ctx *actor.Context) {
	if runner.Status == msg.Status_READY {
		executor := s.executorManager.GetExecutor(runner.Name)
		if executor != nil {
			runner_pid := ctx.Sender()
			senderAddr := runner.Address
			senderPID := actor.NewPID(senderAddr, "runner/"+runner.Name)
			log.Println("Sending deployment object to -> ", runner_pid)
			ctx.Send(senderPID, executor.DeploymentObject)
		}
	} else if runner.Status == msg.Status_COMPLETED {
		executor := s.executorManager.GetExecutor(runner.Name)
		if executor == nil {
			slog.Warn("Executor with name ", runner.Name, " was not found")
			return
		}
		executor.SetStatus(models.SucceededStatus)
		executor.Error = ""
	} else if runner.Status == msg.Status_FAILED {
		executor := s.executorManager.GetExecutor(runner.Name)
		if executor == nil {
			slog.Warn("Executor with name ", runner.Name, " was not found")
			return
		}
		executor.SetStatus(models.FailedStatus)
		executor.Error = runner.GetError()
	} else {
		slog.Warn("Unknown status")
	}
}
