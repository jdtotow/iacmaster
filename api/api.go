package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jdtotow/iacmaster/controllers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port              int
	router            *gin.Engine
	supportedEndpoint []string
	dbController      *controllers.DBController
	seController      *controllers.SecurityController
	system            *controllers.System
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

func CreateServer(port int, dbController *controllers.DBController, seController *controllers.SecurityController, system *controllers.System) *Server {
	return &Server{
		port:              port,
		router:            gin.Default(),
		supportedEndpoint: getSupportedEnpoint(),
		dbController:      dbController,
		seController:      seController,
		system:            system,
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
		s.router.GET(path, s.skittlesMan)           // get all entries
		s.router.GET(path+"/:id", s.skittlesMan)    // get one identify by id
		s.router.DELETE(path+"/:id", s.skittlesMan) // delete one identify by id
		s.router.POST(path, s.skittlesMan)          // create new one
		s.router.PATCH(path+"/:id", s.skittlesMan)  //edit one field of the entry identify by id
		s.router.PUT(path+"/:id", s.skittlesMan)    //replace the entiere object identify by id
	}

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
	method := context.Request.Method

	var objectName string = ""
	for _, _path := range s.supportedEndpoint {
		if strings.HasPrefix(path, _path) {
			objectName = strings.Replace(_path, "/", "", 1)
			s.dbController.Handle(context, objectName)
		}
	}

	context.IndentedJSON(
		http.StatusOK,
		gin.H{
			"path":    path,
			"method":  method,
			"message": "The proper handler will be invoqued",
		},
	)
}
