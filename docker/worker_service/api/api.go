package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"workerservice/controllers"
	"workerservice/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port   int
	router *gin.Engine
	logic  *controllers.Logic
}

func CreateServer() *Server {
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		port = 2020
	}
	return &Server{
		port:   port,
		router: gin.Default(),
		logic:  controllers.CreateLogic("/tmp"),
	}
}

func (s *Server) Start() error {
	url := ":" + fmt.Sprintf("%d", s.port)
	s.router.Use(gin.Recovery())
	s.router.Use(jsonLoggerMiddleware())

	s.router.GET("/", s.homePage)
	s.router.GET("/health", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{"status": "ok"})
	})
	s.router.POST("/deployment", s.addDeployment)
	s.router.POST("/updalod", s.uploadFile)

	log.Println("Starting api server ...")
	err := s.router.Run(url)
	return err
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

func (s *Server) homePage(context *gin.Context) {
	context.IndentedJSON(
		http.StatusOK,
		gin.H{
			"IaCMaster Client API Version": "0.0.1",
		},
	)
}

func (s *Server) addDeployment(context *gin.Context) {
	var deployment *models.Deployment = &models.Deployment{}
	err := context.BindJSON(deployment)
	if err != nil {
		context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	}
	if s.logic.AddDeployment(deployment) {
		context.IndentedJSON(http.StatusCreated, gin.H{"id": deployment.Name})
	} else {
		context.IndentedJSON(http.StatusConflict, gin.H{"error": "Deployment already exists"})
	}
}

func (s *Server) uploadFile(context *gin.Context) {
	form, _ := context.MultipartForm()
	files := form.File["file"]
	environment_id := context.PostForm("environment_id")
	for _, file := range files {
		context.SaveUploadedFile(file, "/tmp/"+environment_id+".tfvars")
	}
}
