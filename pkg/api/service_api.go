package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jdtotow/iacmaster/pkg/controllers"
	"github.com/jdtotow/iacmaster/pkg/models"

	"github.com/gin-gonic/gin"
)

type ServiceServer struct {
	port   int
	router *gin.Engine
	logic  *controllers.Logic
}

func CreateServiceServer() *ServiceServer {
	port, err := strconv.Atoi(os.Getenv("SERVICE_PORT"))
	if err != nil {
		port = 2020
	}
	return &ServiceServer{
		port:   port,
		router: gin.Default(),
		logic:  controllers.CreateLogic("/tmp"),
	}
}

func (s *ServiceServer) Start() error {
	url := ":" + fmt.Sprintf("%d", s.port)
	s.router.Use(gin.Recovery())
	s.router.Use(jsonLoggerMiddleware())

	s.router.GET("/", s.homePage)
	s.router.GET("/health", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{"status": "ok"})
	})
	s.router.POST("/deployment", s.addDeployment)
	s.router.POST("/destroy", s.deleteDeployment)
	s.router.POST("/updalod", s.uploadFile)

	log.Println("Starting api server ...")
	err := s.router.Run(url)
	return err
}

func (s *ServiceServer) homePage(context *gin.Context) {
	context.IndentedJSON(
		http.StatusOK,
		gin.H{
			"IaCMaster Client API Version": "0.0.1",
		},
	)
}

func (s *ServiceServer) addDeployment(context *gin.Context) {
	var deployment *models.Deployment = &models.Deployment{}
	err := context.BindJSON(deployment)
	if err != nil {
		context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	}
	if s.logic.AddDeployment(deployment) {
		context.IndentedJSON(http.StatusCreated, gin.H{"id": deployment.Name})
	} else {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error occured, for details check deployment error"})
	}
}

func (s *ServiceServer) deleteDeployment(context *gin.Context) {
	var deployment *models.Deployment = &models.Deployment{}
	err := context.BindJSON(deployment)
	if err != nil {
		context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	}
	if s.logic.HasDeployment(deployment.Name) {
		s.logic.DeleteDeployment(deployment)
		context.IndentedJSON(http.StatusAccepted, gin.H{})
	} else {
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": "deployment not found"})
	}
}

func (s *ServiceServer) uploadFile(context *gin.Context) {
	form, _ := context.MultipartForm()
	fileType := context.PostForm("fileType")
	files := form.File["file"]
	environment_id := context.PostForm("environment_id")
	filename := ""
	if fileType == "terraform_variables_values" {
		filename = "variables.tfvars"
	} else if fileType == "gcp_credential" {
		filename = "gcp_credential.json"
	} else {
		filename = "send.file"
	}
	for _, file := range files {
		context.SaveUploadedFile(file, "/tmp/"+environment_id+"/"+filename)
	}
}
