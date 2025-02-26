package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jdtotow/iacmaster/controllers"
	"github.com/jdtotow/iacmaster/models"
)

type Server struct {
	port              int
	router            *gin.Engine
	supportedEndpoint []string
	channel           *chan models.HTTPMessage
	dbController      *controllers.DBController
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

func CreateServer(channel *chan models.HTTPMessage) *Server {
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		port = 3000
	}
	return &Server{
		port:              port,
		router:            gin.Default(),
		supportedEndpoint: getSupportedEnpoint(),
		channel:           channel,
		dbController:      controllers.CreateDBController(),
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

func (s *Server) Start() error {
	url := ":" + fmt.Sprintf("%d", s.port)
	s.router.Use(gin.Recovery())
	s.router.Use(jsonLoggerMiddleware())

	s.router.GET("/", s.homePage)
	s.router.POST("/", s.homePage)

	/*
		"/settings",
			"/environment",
			"/organization",
			"/iacartifact",
			"/variable",
			"/cloudcredential",
	*/
	//s.router.POST("/environment/:id/settings", s.createSettings)
	//s.router.GET("/environment/:id/settings/*setting_id", s.getSettings)

	for _, path := range s.supportedEndpoint {
		s.router.GET(path, s.skittlesMan)           // get all entries
		s.router.GET(path+"/:id", s.skittlesMan)    // get one identify by uuid
		s.router.DELETE(path+"/:id", s.skittlesMan) // delete one identify by uuid
		s.router.POST(path, s.skittlesMan)          // create new one
		s.router.PATCH(path+"/:id", s.skittlesMan)  //edit one field of the entry identify by uuid
		s.router.PUT(path+"/:id", s.skittlesMan)    //replace the entiere object identify by uuid
	}
	s.router.POST("/environment/:id/*action", s.deployEnvironment)

	log.Println("Starting api server ...")
	err := s.router.Run(url)
	return err
}

func (s *Server) homePage(context *gin.Context) {
	context.IndentedJSON(
		http.StatusOK,
		gin.H{
			"IaCMaster API Version": "0.0.1",
		},
	)
}

func (s *Server) skittlesMan(context *gin.Context) {
	path := context.FullPath()
	var objectName string = ""
	for _, _path := range s.supportedEndpoint {
		if strings.HasPrefix(path, _path) {
			objectName = strings.Replace(_path, "/", "", 1)
			token, _ := context.Cookie("Authorization")
			metadata := map[string]string{
				"action":    "no_action",
				"object_id": "no",
			}
			s.Handle(context, objectName)
			message := models.HTTPMessage{
				ObjectName:    objectName,
				RequestOrigin: context.ClientIP(),
				Method:        context.Request.Method,
				Url:           context.Request.RequestURI,
				Token:         token,
				Body:          context.Request.Body,
				Params:        context.Request.URL.Query(),
				Metadata:      metadata,
			}
			*s.channel <- message
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{})
}

func (s *Server) Handle(context *gin.Context, objectName string) {
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

func (s *Server) deployEnvironment(context *gin.Context) {
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
		metadata := map[string]string{
			"action":    "create_env",
			"object_id": id,
		}
		message := models.HTTPMessage{
			ObjectName:    "environment",
			RequestOrigin: context.ClientIP(),
			Method:        context.Request.Method,
			Url:           context.Request.RequestURI,
			Token:         "-",
			Body:          context.Request.Body,
			Params:        context.Request.URL.Query(),
			Metadata:      metadata,
		}
		*s.channel <- message
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

	} else {

	}
}
