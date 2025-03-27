package api

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/gin-gonic/gin"
	"github.com/jdtotow/iacmaster/pkg/controllers"
	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type SystemServer struct {
	port              int
	router            *gin.Engine
	supportedEndpoint []string
	dbController      *controllers.DBController
	systemPID         string
	systemAddr        string
	actorEngine       *actor.Engine
	nodeInfo          *msg.NodeInfo
}

func getSupportedEnpoint() []string {
	return []string{
		//"/",
		"/user",
		"/group",
		"/project",
		"/organization",
		"/iacartifact",
		"role",
		"/token",
		"/cloudcredential",
		"/environment",
		"/settings",
	}
}

func CreateSystemServer() *SystemServer {
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		port = 3000
	}
	return &SystemServer{
		port:              port,
		router:            gin.Default(),
		supportedEndpoint: getSupportedEnpoint(),
		dbController:      controllers.CreateDBController(),
		systemPID:         "iacmaster/system",
		systemAddr:        os.Getenv("IACMASTER_SYSTEM_ADDRESS") + ":" + os.Getenv("IACMASTER_SYSTEM_PORT"),
	}
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func (s *SystemServer) Start() *SystemServer {
	url := ":" + fmt.Sprintf("%d", s.port)
	s.router.Use(gin.Recovery())
	s.router.Use(jsonLoggerMiddleware())

	s.router.GET("/", s.homePage)
	s.router.POST("/", s.homePage)
	s.router.GET("/nodetype", s.getNodeType)

	for _, path := range s.supportedEndpoint {
		s.router.GET(path, s.skittlesMan)           // get all entries
		s.router.GET(path+"/:id", s.skittlesMan)    // get one identify by uuid
		s.router.DELETE(path+"/:id", s.skittlesMan) // delete one identify by uuid
		s.router.POST(path, s.skittlesMan)          // create new one
		s.router.PATCH(path+"/:id", s.skittlesMan)  //edit one field of the entry identify by uuid
		s.router.PUT(path+"/:id", s.skittlesMan)    //replace the entiere object identify by uuid
	}
	s.router.POST("/environment/:id/*action", s.deployEnvironment)

	log.Println("Starting api System Server ...")
	err := s.router.Run(url)
	if err != nil {
		return nil
	}
	return s
}

func (s *SystemServer) homePage(context *gin.Context) {
	context.IndentedJSON(
		http.StatusOK,
		gin.H{
			"IaCMaster API Version": "0.0.1",
		},
	)
}

func (s *SystemServer) skittlesMan(context *gin.Context) {
	path := context.FullPath()
	var objectName string = ""
	for _, _path := range s.supportedEndpoint {
		if strings.HasPrefix(path, _path) {
			objectName = strings.Replace(_path, "/", "", 1)
			s.Handle(context, objectName)
			var operation msg.Operation
			operation.Action = "no_action"
			operation.ObjectID = "no"
			systemPID := actor.NewPID(s.systemAddr, s.systemPID)
			s.actorEngine.Send(systemPID, &operation)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{})
}

func (s *SystemServer) Handle(context *gin.Context, objectName string) {
	if context.Request.Method == "POST" {
		if objectName == "organization" {
			//create an organization
			var org *models.Organization = &models.Organization{}
			err := context.BindJSON(org)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(org)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{"id": org.ID})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else if objectName == "user" {
			var user *models.User = &models.User{}
			err := context.BindJSON(user)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(user)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else if objectName == "project" {
			var project *models.Project = &models.Project{}
			err := context.BindJSON(project)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(project)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{"id": project.ID})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else if objectName == "token" {
			var token *models.Token = &models.Token{}
			err := context.BindJSON(token)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(token)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{"id": token.ID})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else if objectName == "iacartifact" {
			var arti *models.IaCArtifact = &models.IaCArtifact{}
			err := context.BindJSON(arti)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(arti)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{"id": arti.ID})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else if objectName == "environment" {
			var env *models.Environment = &models.Environment{}
			err := context.BindJSON(env)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(env)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{"id": env.ID})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else if objectName == "settings" {
			var settings *models.IaCExecutionSettings = &models.IaCExecutionSettings{}
			err := context.BindJSON(settings)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(settings)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{"id": settings.ID})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else if objectName == "cloudcredential" {
			var credential *models.CloudCredential = &models.CloudCredential{}
			err := context.BindJSON(credential)
			if err != nil {
				context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			}
			result := s.dbController.CreateInstance(credential)
			if result.Error == nil {
				context.IndentedJSON(http.StatusCreated, gin.H{"id": credential.ID})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": "object handler not found"})
		}
	} else if context.Request.Method == "GET" {
		id := context.Param("id")
		if objectName == "organization" {
			if id == "" {
				orgs := []models.Organization{}
				result := s.dbController.GetAll(&orgs)
				if result.Error != nil {
					context.IndentedJSON(http.StatusNotFound, gin.H{})
					return
				}
				context.JSON(http.StatusOK, orgs)
			} else {
				org := models.Organization{}
				result := s.dbController.GetObjectByID(&org, id)
				if result.Error != nil {
					context.IndentedJSON(http.StatusNotFound, gin.H{})
					return
				}
				context.JSON(http.StatusOK, org)
			}
		} else if objectName == "user" {
			if id == "" {
				users := []models.User{}
				result := s.dbController.GetAll(&users)
				if result.Error != nil {
					context.IndentedJSON(http.StatusNotFound, gin.H{})
					return
				}
				context.JSON(http.StatusOK, users)
			} else {
				user := models.User{}
				result := s.dbController.GetObjectByID(&user, id)
				if result.Error != nil {
					context.IndentedJSON(http.StatusNotFound, gin.H{})
					return
				}
				context.JSON(http.StatusOK, user)
			}
		} else if objectName == "project" {
			if id == "" {
				projects := []models.Project{}
				result := s.dbController.GetAll(&projects)
				if result.Error != nil {
					context.IndentedJSON(http.StatusNotFound, gin.H{})
					return
				}
				context.JSON(http.StatusOK, projects)
			} else {
				project := models.Project{}
				result := s.dbController.GetObjectByID(&project, id)
				if result.Error != nil {
					context.IndentedJSON(http.StatusNotFound, gin.H{})
					return
				}
				context.JSON(http.StatusOK, project)
			}
		} else {
			context.IndentedJSON(http.StatusNotFound, gin.H{})
		}

	} else if context.Request.Method == "PUT" {

	} else if context.Request.Method == "PATCH" {

	} else if context.Request.Method == "DELETE" {

	} else {
		context.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
	}
}

func (s *SystemServer) deployEnvironment(context *gin.Context) {
	id := context.Param("id")
	action := context.Param("action")
	env := models.Environment{}
	result := s.dbController.GetObjectByID(&env, id)
	if result.Error != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}
	if action == "/deploy" {
		log.Println("Deployment of the environment with ID", id)
		var operation msg.Operation
		operation.Action = "create_env"
		operation.ObjectID = id
		systemPID := actor.NewPID(s.systemAddr, s.systemPID)
		s.actorEngine.Send(systemPID, &operation)

	} else if action == "/variables" {
		form, _ := context.MultipartForm()
		files := form.File["file"]
		pwd, _ := os.Getwd()
		environment_id := context.PostForm("environment_id")
		for _, file := range files {
			log.Println(file.Filename)
			// Upload the file to specific dst.
			context.SaveUploadedFile(file, pwd+"/tmp/"+environment_id+".tfvars")
		}

	} else if action == "/destroy" {
		log.Println("Destroying of the environment with ID", id)
		var operation msg.Operation
		operation.Action = "destroy_env"
		operation.ObjectID = id
		systemPID := actor.NewPID(s.systemAddr, s.systemPID)
		s.actorEngine.Send(systemPID, &operation)

	} else {

	}
}

func (s *SystemServer) getNodeType(context *gin.Context) {
	nodeType := "primary"
	if s.nodeInfo.NodeType == uint32(models.Secondary) {
		nodeType = "secondary"
	}
	context.IndentedJSON(http.StatusOK, nodeType)
}

func (s *SystemServer) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		log.Println("API actor started at -> ", ctx.Engine().Address())
		s.actorEngine = ctx.Engine()
		go s.Start()
	case actor.Stopped:
		log.Println("API actor has stopped")
	case actor.Initialized:
		log.Println("API actor initialized")
	case *actor.PID:
		log.Println("API actor has god an ID")
	case *msg.NodeInfo:
		s.nodeInfo = m
		log.Println("API has received node info message received -> ", m)
	default:
		slog.Warn("API Server got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}
