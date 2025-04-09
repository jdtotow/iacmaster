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
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	nodeAttributes    *msg.NodeAttribute
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
	s.router.Use(cors.Default())

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
	s.router.POST("/environment/:id/*action", s.environmentActions)

	log.Println("Starting api System Server ...")
	err := s.router.Run(url)
	if err != nil {
		return nil
	}
	return s
}

func (s *SystemServer) CreateEmptyEntityInstance(objectName string) interface{} {
	if objectName == "organization" {
		return &models.Organization{}
	} else if objectName == "user" {
		return &models.User{}
	} else if objectName == "project" {
		return &models.Project{}
	} else if objectName == "token" {
		return &models.Token{}
	} else if objectName == "iacartifact" {
		return &models.IaCArtifact{}
	} else if objectName == "environment" {
		return &models.Environment{}
	} else if objectName == "settings" {
		return &models.IaCExecutionSettings{}
	} else if objectName == "cloudcredential" {
		return &models.CloudCredential{}
	} else {
		return nil
	}
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

func (s *SystemServer) VerifyJWT(tokenStr string) (*jwt.MapClaims, error) {
	jwtSecret := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s *SystemServer) CreateJWTToken(userObject any) (string, error) {

	user := userObject.(*models.User)
	expirationTime := time.Now().Add(24 * time.Hour)
	jwtSecret := []byte(os.Getenv("SECRET_KEY"))
	// Create the claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "iacmaster",
	}

	// Create the token using the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *SystemServer) Handle(context *gin.Context, objectName string) {
	object := s.CreateEmptyEntityInstance(objectName)
	if object == nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}

	if context.Request.Method == "GET" {
		id := context.Param("id")
		if id == "" {
			result := s.dbController.GetAll(&object)
			if result.Error != nil {
				context.IndentedJSON(http.StatusNotFound, gin.H{})
				return
			}
			context.JSON(http.StatusOK, object)
		} else {
			result := s.dbController.GetObjectByID(&object, id)
			if result.Error != nil {
				context.IndentedJSON(http.StatusNotFound, gin.H{})
				return
			}
			context.JSON(http.StatusOK, object)
		}
	} else if context.Request.Method == "POST" {
		err := context.BindJSON(object)
		if err != nil {
			context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		}
		result := s.dbController.CreateInstance(object)
		if result.Error == nil {
			if objectName == "user" {
				token, err := s.CreateJWTToken(object)
				if err != nil {
					context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				} else {
					context.IndentedJSON(http.StatusCreated, gin.H{"token": token, "object": object})
					return
				}
			}
			context.IndentedJSON(http.StatusCreated, gin.H{"object": object})
		} else {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		}
	} else if context.Request.Method == "DELETE" {
		id := context.Param("id")
		if id == "" {
			context.IndentedJSON(http.StatusNotFound, gin.H{})
			return
		}
		result := s.dbController.GetClient().Delete(object, "ID = ?", id)
		if result.Error == nil {
			context.IndentedJSON(http.StatusOK, gin.H{"object": object})
		} else {
			context.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
		}
	} else if context.Request.Method == "PATCH" {
		id := context.Param("id")
		if id == "" {
			context.IndentedJSON(http.StatusNotFound, gin.H{})
			return
		}
		existing_object := s.CreateEmptyEntityInstance(objectName)
		result := s.dbController.GetObjectByID(&existing_object, id)
		if result.Error != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{})
			return
		}
		err := context.BindJSON(object)
		if err != nil {
			context.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		}
		result = s.dbController.GetClient().Model(&existing_object).Updates(object)
		if result.Error != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			context.IndentedJSON(http.StatusOK, gin.H{})
		}

	} else {
		context.IndentedJSON(http.StatusMethodNotAllowed, gin.H{})
	}
}

func (s *SystemServer) environmentActions(context *gin.Context) {
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
	case *msg.NodeAttribute:
		s.nodeAttributes = m
		log.Println("API has received node attributes message received -> name: ", m.NodeName, ", attribute: ", m.Attribute)
	default:
		slog.Warn("API Server got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}
