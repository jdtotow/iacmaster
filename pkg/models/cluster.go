package models

import (
	"fmt"
	"math/rand/v2"

	"slices"

	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type ClusterInfo struct {
	Name  string
	Nodes []*msg.NodeInfo
}

func CreateClusterInfo(name string) *ClusterInfo {
	return &ClusterInfo{
		Name:  name,
		Nodes: []*msg.NodeInfo{},
	}
}

func (ci *ClusterInfo) NodeExist(name string) bool {
	for _, node := range ci.Nodes {
		if node.Name == name {
			return true
		}
	}
	return false
}

func (ci *ClusterInfo) AddNode(node *msg.NodeInfo) {
	if ci.NodeExist(node.Name) {
		return
	}
	ci.Nodes = append(ci.Nodes, node)
}

func (ci *ClusterInfo) RemoveNode(name string) {
	for i, node := range ci.Nodes {
		if node.Name == name {
			ci.Nodes = append(ci.Nodes[:i], ci.Nodes[i+1:]...)
			break
		}
	}
}

func (ci *ClusterInfo) GetNodeExecutingDeployment(deployment string) *msg.NodeInfo {
	for _, node := range ci.Nodes {
		if slices.Contains(node.Deployments, deployment) {
			return node
		}
	}
	return nil
}

func (ci *ClusterInfo) AddDeploymentToNode(nodeName, deployment string) bool {
	for _, node := range ci.Nodes {
		if node.Name == nodeName {
			node.Deployments = append(node.Deployments, deployment)
			return true
		}
	}
	return false
}
func (ci *ClusterInfo) GetNodes() []*msg.NodeInfo {
	return ci.Nodes
}

func (ci *ClusterInfo) GetRandomNode() *msg.NodeInfo {
	rand_int := rand.IntN(len(ci.Nodes)-0) + 0
	return ci.Nodes[rand_int]
}

func (ci *ClusterInfo) GetNodeByName(name string) *msg.NodeInfo {
	for _, node := range ci.Nodes {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func (ci *ClusterInfo) GetNodeAddr(name string) (string, error) {
	for _, node := range ci.Nodes {
		if node.Name == name {
			return node.Addr, nil
		}
	}
	return "", fmt.Errorf("node not found")
}

func (ci *ClusterInfo) GetPIDNode(nodeName, process string) *actor.PID {
	for _, node := range ci.Nodes {
		if node.Name == nodeName {
			return actor.NewPID(node.Addr, "iacmaster/"+process)
		}
	}
	return nil
}
