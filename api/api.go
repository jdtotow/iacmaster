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
		"/settings",
		"/environment",
		"/organization",
		"/iacartifact",
		"/variable",
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

	for _, path := range s.supportedEndpoint {
		s.router.GET(path, s.skittlesMan)             // get all entries
		s.router.GET(path+"/:uuid", s.skittlesMan)    // get one identify by uuid
		s.router.DELETE(path+"/:uuid", s.skittlesMan) // delete one identify by uuid
		s.router.POST(path, s.skittlesMan)            // create new one
		s.router.PATCH(path+"/:uuid", s.skittlesMan)  //edit one field of the entry identify by uuid
		s.router.PUT(path+"/:uuid", s.skittlesMan)    //replace the entiere object identify by uuid
	}
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
			s.Handle(context, objectName)
			message := models.HTTPMessage{
				ObjectName:    objectName,
				RequestOrigin: context.ClientIP(),
				Method:        context.Request.Method,
				Url:           context.Request.RequestURI,
				Token:         token,
				Body:          context.Request.Body,
				Params:        context.Request.URL.Query(),
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
				context.IndentedJSON(http.StatusCreated, gin.H{})
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
				context.IndentedJSON(http.StatusCreated, gin.H{})
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
				context.IndentedJSON(http.StatusCreated, gin.H{})
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
				context.IndentedJSON(http.StatusCreated, gin.H{})
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
				context.IndentedJSON(http.StatusCreated, gin.H{})
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}
		} else {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": "object handler not found"})
		}
	} else if context.Request.Method == "GET" {

	} else if context.Request.Method == "PUT" {

	} else if context.Request.Method == "PATCH" {

	} else if context.Request.Method == "DELETE" {

	} else {
		context.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
	}
}
