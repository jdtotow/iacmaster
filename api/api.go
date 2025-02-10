package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port              int
	router            *gin.Engine
	supportedEndpoint []string
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

func CreateServer(port int) *Server {
	return &Server{
		port:              port,
		router:            gin.Default(),
		supportedEndpoint: getSupportedEnpoint(),
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
		s.router.GET(path, s.skittlesMan)
		s.router.POST(path, s.skittlesMan)
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

func (s *Server) isPathMatch(path, fullpath string) bool {
	if !strings.HasPrefix(fullpath, path) {
		return false
	}
	if len(fullpath) == len(path) {
		return true
	}
	if fullpath[len(path)] == '/' {
		return true
	}
	return false
}

func (s *Server) skittlesMan(context *gin.Context) {
	path := context.FullPath()
	method := context.Request.Method

	for _, _path := range s.supportedEndpoint {
		if s.isPathMatch(_path, path) {
			context.IndentedJSON(
				http.StatusOK,
				gin.H{
					"path":    path,
					"matched": _path,
					"method":  method,
					"message": "The proper handler will be invoqued",
				},
			)
			return
		}
	}

	context.IndentedJSON(
		http.StatusNotFound,
		gin.H{
			"error": "path or method not supported",
		},
	)
}
