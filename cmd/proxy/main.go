package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/madflojo/tasks"
)

type NodeType uint32

const (
	Primary NodeType = iota + 1
	Secondary
	Unknown
)

type Node struct {
	Name     string
	Addr     string
	NodeType NodeType
}

func (n *Node) SetType(_type NodeType) {
	n.NodeType = _type
}

func NewNode(name, addr string, _type NodeType) *Node {
	return &Node{
		Name:     name,
		Addr:     addr,
		NodeType: _type,
	}
}

func createReverseProxy(target string) http.Handler {
	// Parse the target URL
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}
	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// Customize the request
	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", targetURL.Host)
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = targetURL.Path + req.URL.Path
	}

	// Customize the response
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Add("X-Proxy-Server", "IaCMaster-Reverse-Proxy")
		return nil
	}

	return proxy
}

type ReverseProxy struct {
	Port      string
	Nodes     []*Node
	router    *gin.Engine
	scheduler *tasks.Scheduler
}

func (p *ReverseProxy) NodeSize() int {
	return len(p.Nodes)
}

func (p *ReverseProxy) GetPrimaryNode() *Node {
	for _, node := range p.Nodes {
		if node.NodeType == Primary {
			return node
		}
	}
	return nil
}

func GetNodes() []*Node {
	var result []*Node
	for _, chunk := range strings.Split(os.Getenv("CLUSTER"), ",") {
		setting := strings.Split(chunk, "=")
		nodeName := setting[0]
		nodeAddr := setting[1] + ":3000"
		node := NewNode(nodeName, nodeAddr, Secondary)
		result = append(result, node)
	}
	return result
}
func NewReverseProxy(port string) *ReverseProxy {
	return &ReverseProxy{
		Port:      port,
		router:    gin.Default(),
		Nodes:     GetNodes(),
		scheduler: tasks.New(),
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

func (p *ReverseProxy) UpdatePrimary() error {
	client := &http.Client{Timeout: 2 * time.Second} // Set timeout to 2 seconds per request

	for _, node := range p.Nodes {
		url := fmt.Sprintf("http://%s/nodetype", node.Addr)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error requesting %s: %v\n", url, err)
			continue // Skip this host if there's an error
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response from %s: %v\n", url, err)
			continue
		}
		var nodeType NodeType
		body_str := strings.ReplaceAll(string(body), "\"", "")
		if body_str == "primary" {
			nodeType = Primary
		} else if string(body) == "secondary" {
			nodeType = Secondary
		} else {
			nodeType = Unknown
		}
		node.SetType(nodeType)
	}
	return nil
}

func (p *ReverseProxy) callProxy(context *gin.Context) {
	primary_node := p.GetPrimaryNode()
	if primary_node == nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": "no primary node found, please try again later"})
		return
	}
	fmt.Println("Forwarding request to -> ", primary_node.Addr)
	proxy := createReverseProxy("http://" + primary_node.Addr)
	proxy.ServeHTTP(context.Writer, context.Request)
}

func (p *ReverseProxy) Start() {
	url := ":" + p.Port
	p.router.Use(gin.Recovery())
	p.router.Use(jsonLoggerMiddleware())
	p.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Your Vue.js dev server address
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	p.router.Any("/*any", p.callProxy)

	log.Println("Starting api System Server ...")

	_ = p.scheduler.AddWithID("UpdateNodeType", &tasks.Task{
		Interval: time.Duration(30 * time.Second),
		TaskFunc: p.UpdatePrimary,
		ErrFunc:  nil,
	})

	err := p.router.Run(url)
	if err != nil {
		log.Fatal("Could not start reverse proxy server : ", err.Error())
	}
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Could not load .env")
	}
}

func main() {
	port := os.Getenv("REVERSE_PROXY_SERVER_PORT")
	if port == "" {
		port = "5454"
	}

	proxy := NewReverseProxy(port)
	proxy.Start()
}
