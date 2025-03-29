package models

import (
	"errors"
	"log"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/rpc"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
	"github.com/madflojo/tasks"
)

type NodeType uint32

const (
	Primary NodeType = iota + 1
	Secondary
)

type NodeStatus uint32

const (
	Init NodeStatus = iota + 1
	Ready
	Running
	Down
	Failed
	Unknown
)

type NodeAttribute uint32

const (
	ExecutorNodeAttribute NodeAttribute = iota + 1
	ManagerNodeAttribute
	LoggingNodeAttribute
	APINodeAttribute
)

type NodeMode uint32

const (
	Standalone NodeMode = iota + 1
	Cluster
)
const (
	ScheduleHeartbeat              = "heartbeat"
	ScheduleLeaderHeartbeatTimeout = "leader_heartbeat_timeout"
	CheckPeers                     = "check_peers"
)

type Node struct {
	Type                   NodeType
	Name                   string
	Status                 NodeStatus
	Mode                   NodeMode
	Addr                   string
	Peers                  *Peers
	Attributes             []NodeAttribute
	MaxRetry               int
	ActorEngine            *actor.Engine
	scheduler              *tasks.Scheduler
	LeaderHeartbeatTimeout time.Duration
	HeartbeatInterval      time.Duration
	mtex                   sync.Mutex
}

func GenerateLeaderHeartbeatTimeout() time.Duration {
	min := 7
	max := 20
	rand_int := rand.IntN(max-min) + min
	return time.Duration(rand_int * int(time.Second))
}

func NewNode() *Node {
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		log.Fatal("node name must be set, ex.: NODE_NAME=node-01")
	}
	addr := os.Getenv("IACMASTER_SYSTEM_ADDRESS") + ":" + os.Getenv("IACMASTER_NODE_PORT")
	var mode NodeMode
	if os.Getenv("CLUSTER") != "" {
		mode = Cluster
	} else {
		mode = Standalone
	}

	return &Node{
		Type:                   Secondary, //2 as secondary node
		Name:                   nodeName,
		Mode:                   mode,
		Status:                 Init,
		Addr:                   addr,
		Attributes:             []NodeAttribute{},
		Peers:                  NewPeers(),
		MaxRetry:               3,
		LeaderHeartbeatTimeout: GenerateLeaderHeartbeatTimeout(),
		HeartbeatInterval:      time.Duration(5 * time.Second),
		scheduler:              tasks.New(),
	}
}
func (n *Node) GetType() NodeType {
	return n.Type
}
func (n *Node) GetName() string {
	return n.Name
}
func (n *Node) GetNodeStatus() NodeStatus {
	return n.Status
}
func (n *Node) SetType(_type NodeType) {
	n.Type = _type
	type_str := "secondary"
	if _type == Primary {
		type_str = "primary"
	}
	log.Println("Node : ", n.Name, " has become : ", type_str)
}
func (n *Node) SetName(name string) {
	n.Name = name
}
func (n *Node) SetNodeStatus(status NodeStatus) {
	n.Status = status
}
func (n *Node) SetNodeMode(mode NodeMode) {
	n.Mode = mode
}

func (node *Node) NewListener() (net.Listener, error) {
	addr, err := net.Listen("tcp", node.Addr)
	return addr, err
}

func (node *Node) ConnectToPeers() {
	for _, chunk := range strings.Split(os.Getenv("CLUSTER"), ",") {
		setting := strings.Split(chunk, "=")
		peerID := setting[0]
		peerAddr := setting[1] + ":" + os.Getenv("IACMASTER_NODE_PORT")
		if node.IsItself(peerID) {
			continue //skipping myself
		}
		rpcClient := node.connect(peerAddr)
		if rpcClient == nil {
			continue
		}
		pingNodeMessage := NodeMessage{FromPeerID: node.Name, Type: PING}
		reply, _ := node.CommunicateWithPeer(rpcClient, pingNodeMessage)

		if reply.IsPongMessage() {
			log.Printf("%s got pong NodeMessage from %s", node.Name, peerID)
			node.Peers.Add(peerID, rpcClient)
		}

	}
}

func (node *Node) connect(peerAddr string) *rpc.Client {
	for range node.MaxRetry {
		client, err := rpc.Dial("tcp", peerAddr)
		if err != nil {
			log.Printf("Error dialing rpc dial %s", err.Error())
			time.Sleep(50 * time.Millisecond)
			continue
		}
		return client
	}
	return nil
}

func (node *Node) CommunicateWithPeer(RPCClient *rpc.Client, args NodeMessage) (NodeMessage, error) {
	var reply NodeMessage

	err := RPCClient.Call("Node.HandleNodeMessage", args, &reply)
	if err != nil {
		log.Printf("Error calling HandleNodeMessage %s", err.Error())
	}

	return reply, err
}

func (node *Node) ProcessLeaderHeartbeat() {
	node.scheduler.Del(ScheduleLeaderHeartbeatTimeout)
	_ = node.scheduler.AddWithID(ScheduleLeaderHeartbeatTimeout, &tasks.Task{
		Interval: node.LeaderHeartbeatTimeout,
		TaskFunc: node.leaderHeartbeatTimeoutFunc,
		ErrFunc:  node.leaderHeartbeatTimeoutErrorFunc,
	})
}

func (node *Node) HandleNodeMessage(args NodeMessage, reply *NodeMessage) error {
	reply.FromPeerID = node.Name

	switch args.Type {
	case ELECTION:
		reply.Type = ALIVE
	case ELECTED:
		leaderID := args.FromPeerID
		log.Printf("Election is done. %s has a new leader %s", node.Name, leaderID)
		reply.Type = OK
	case PING:
		reply.Type = PONG
	case LEADERHEARTBEAT:
		node.ProcessLeaderHeartbeat()
		reply.Type = OK
	}
	return nil
}

func (node *Node) CheckPeers() error {
	peers := node.Peers.ToList()
	if len(peers) < 1 {
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(len(peers))
	for i := range peers {
		peer := peers[i]
		go func(peer Peer) {
			pingMessage := NodeMessage{FromPeerID: node.Name, Type: PING}
			reply, _ := node.CommunicateWithPeer(peer.RPCClient, pingMessage)
			if !reply.IsPongMessage() {
				log.Println("Peer : ", peer.ID, " is not responding, it will be deleted")
				node.Peers.Delete(peer.ID)
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()
	return nil
}

func (node *Node) Elect() error {
	isHighestRankedNodeAvailable := false

	peers := node.Peers.ToList()
	for i := range peers {
		peer := peers[i]

		if node.IsRankHigherThan(peer.ID) {
			continue
		}

		log.Printf("%s send ELECTION NodeMessage to peer %s", node.Name, peer.ID)
		electionNodeMessage := NodeMessage{FromPeerID: node.Name, Type: ELECTION}

		reply, _ := node.CommunicateWithPeer(peer.RPCClient, electionNodeMessage)

		if reply.IsAliveMessage() {
			isHighestRankedNodeAvailable = true
			node.SetType(Secondary)
			task, _ := node.scheduler.Lookup(CheckPeers)
			if task != nil {
				node.scheduler.Del(CheckPeers)
			}
		}
	}

	if !isHighestRankedNodeAvailable {
		node.SetType(Primary)
		leaderID := node.Name
		electedNodeMessage := NodeMessage{FromPeerID: leaderID, Type: ELECTED}
		node.BroadcastNodeMessage(electedNodeMessage)
		log.Printf("%s is a new leader", node.Name)
		//
		_ = node.scheduler.AddWithID(ScheduleHeartbeat, &tasks.Task{
			Interval: node.HeartbeatInterval,
			TaskFunc: node.heartbeatFunc,
			ErrFunc:  node.heartbeatErrorFunc,
		})
		//
		node.scheduler.Del(ScheduleLeaderHeartbeatTimeout)
		_ = node.scheduler.AddWithID(CheckPeers, &tasks.Task{
			Interval: time.Duration(5 * time.Second),
			TaskFunc: node.CheckPeers,
			ErrFunc:  nil,
		})
		//Sending message to the API server
	}
	node.SendNodeInfoMessage(Running)
	return nil
}

func (node *Node) BroadcastNodeMessage(args NodeMessage) ([]*NodeMessage, error) {
	peers := node.Peers.ToList()
	responses := []*NodeMessage{}
	for i := range peers {
		peer := peers[i]
		reply, _ := node.CommunicateWithPeer(peer.RPCClient, args)
		if reply.FromPeerID != "" {
			responses = append(responses, &reply)
		}
	}
	if len(peers) != len(responses) {
		return []*NodeMessage{}, errors.New("an error occured")
	}
	return responses, nil
}

func (node *Node) IsRankHigherThan(id string) bool {
	return strings.Compare(node.Name, id) == 1
}

func (node *Node) IsItself(id string) bool {
	return node.Name == id
}

func (node *Node) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		log.Println("Node actor started at -> ", ctx.Engine().Address())
		node.ActorEngine = ctx.Engine()
		node.Start()
	case actor.Stopped:
		log.Println("Node actor has stopped")
	case actor.Initialized:
		log.Println("Node actor initialized")
	case *actor.PID:
		log.Println("Node actor has god an ID")
	default:
		slog.Warn("Node got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}

func (node *Node) Start() {
	//node.EventBus.Subscribe(event.LeaderElected, node.PingLeaderContinuously)
	listener, err := node.NewListener()
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	rpcServer := rpc.NewServer()
	rpcServer.Register(node)
	go rpcServer.Accept(listener)
	node.ConnectToPeers()
	log.Printf("%s is aware of own peers %s", node.Name, node.Peers.ToIDs())
	warmupTime := 5 * time.Second
	time.Sleep(warmupTime)
	//node.Elect()
	node.SetNodeStatus(Ready)
	//sending status message to system
	node.SendNodeInfoMessage(Ready)

	log.Println("HeartbeatTimeout is scheduled")
	_ = node.scheduler.AddWithID(ScheduleLeaderHeartbeatTimeout, &tasks.Task{
		Interval: node.LeaderHeartbeatTimeout,
		TaskFunc: node.leaderHeartbeatTimeoutFunc,
		ErrFunc:  node.leaderHeartbeatTimeoutErrorFunc,
	})
}

func (node *Node) SendNodeInfoMessage(status NodeStatus) {
	msg := &msg.NodeInfo{
		Name:       node.Name,
		NodeStatus: uint32(status),
		NodeType:   uint32(node.Type),
		Addr:       os.Getenv("IACMASTER_SYSTEM_ADDRESS"),
	}
	senderAddr := os.Getenv("IACMASTER_SYSTEM_ADDRESS") + ":" + os.Getenv("IACMASTER_SYSTEM_PORT")
	systemPID := actor.NewPID(senderAddr, "iacmaster/system")
	apiPID := actor.NewPID(senderAddr, "iacmaster/api")

	node.ActorEngine.Send(systemPID, msg)
	node.ActorEngine.Send(apiPID, msg)
}

// when no heartbeats are received from leader it triggers a leader election in new term.
// node declares itself as a candidate and sends out proposals to its peers.
func (node *Node) leaderHeartbeatTimeoutFunc() error {
	log.Println("Heartbeat Timeout is called")
	return node.Elect()
}

// if candidate proposal is not successful, the node continues to be a follower.
func (node *Node) leaderHeartbeatTimeoutErrorFunc(err error) {
	node.mtex.Lock()
	defer node.mtex.Unlock()
	node.Type = Secondary
	slog.Error("error encountered on task leader_heartbeat_timeout", "error", err)
}

// triggering of heartbeat task will result in heartbeat signals
// to fellow peers in the cluster.
func (node *Node) heartbeatFunc() error {
	electedNodeMessage := NodeMessage{FromPeerID: node.Name, Type: LEADERHEARTBEAT}
	node.BroadcastNodeMessage(electedNodeMessage)
	return nil
}

// If leader encounters errors when sending heartbeats
// the node will be downgraded to be a follower.
func (node *Node) heartbeatErrorFunc(err error) {
	node.scheduler.Del(ScheduleHeartbeat)

	slog.Info("reinstating leader heartbeat timeout scheduled task")
	_ = node.scheduler.AddWithID(ScheduleLeaderHeartbeatTimeout, &tasks.Task{
		Interval: node.LeaderHeartbeatTimeout,
		TaskFunc: node.leaderHeartbeatTimeoutFunc,
		ErrFunc:  node.leaderHeartbeatTimeoutErrorFunc,
	})
}
