package models

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/jdtotow/iacmaster/pkg/event"
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
)

type NodeMode uint32

const (
	Standalone NodeMode = iota + 1
	Cluster
)

type Node struct {
	Type       NodeType
	Name       string
	Status     NodeStatus
	Mode       NodeMode
	Addr       string
	Peers      *Peers
	EventBus   event.Bus
	Attributes []NodeAttribute
	MaxRetry   int
}

func (n Node) GetType() NodeType {
	return n.Type
}
func (n Node) GetName() string {
	return n.Name
}
func (n Node) GetNodeStatus() NodeStatus {
	return n.Status
}
func (n *Node) SetType(_type NodeType) {
	n.Type = _type
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

func (node *Node) HandleNodeMessage(args NodeMessage, reply *NodeMessage) error {
	reply.FromPeerID = node.Name

	switch args.Type {
	case ELECTION:
		reply.Type = ALIVE
	case ELECTED:
		leaderID := args.FromPeerID
		log.Printf("Election is done. %s has a new leader %s", node.Name, leaderID)
		node.EventBus.Emit(event.LeaderElected, leaderID)
		reply.Type = OK
	case PING:
		reply.Type = PONG
	}

	return nil
}

func (node *Node) Elect() {
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
			node.Type = Secondary
		}
	}

	if !isHighestRankedNodeAvailable {
		leaderID := node.Name
		electedNodeMessage := NodeMessage{FromPeerID: leaderID, Type: ELECTED}
		node.BroadcastNodeMessage(electedNodeMessage)
		log.Printf("%s is a new leader", node.Name)
		node.Type = Primary
	}
}

func (node *Node) BroadcastNodeMessage(args NodeMessage) {
	peers := node.Peers.ToList()
	for i := range peers {
		peer := peers[i]
		node.CommunicateWithPeer(peer.RPCClient, args)
	}
}

func (node *Node) PingLeaderContinuously(_ string, payload any) {
	leaderID := payload.(string)

ping:
	leader := node.Peers.Get(leaderID)
	if leader == nil {
		log.Printf("%s, %s, %s", node.Name, leaderID, node.Peers.ToIDs())
		return
	}

	pingNodeMessage := NodeMessage{FromPeerID: node.Name, Type: PING}
	reply, err := node.CommunicateWithPeer(leader.RPCClient, pingNodeMessage)
	if err != nil {
		log.Println("Leader is down, new election about to start!")
		node.Peers.Delete(leaderID)
		node.Elect()
		return
	}

	if reply.IsPongMessage() {
		log.Printf("Leader %s sent PONG NodeMessage", reply.FromPeerID)
		time.Sleep(3 * time.Second)
		goto ping
	}
}

func (node *Node) IsRankHigherThan(id string) bool {
	return strings.Compare(node.Name, id) == 1
}

func (node *Node) IsItself(id string) bool {
	return node.Name == id
}
