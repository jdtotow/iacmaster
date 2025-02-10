package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port   int
	router *gin.Engine
}

func CreateServer(port int) *Server {
	return &Server{
		port:   port,
		router: gin.Default(),
	}
}
func (s *Server) Start() error {
	url := ":" + fmt.Sprintf("%s", s.port)
	s.router.Run(url)
}
