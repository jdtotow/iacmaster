package models

import "gorm.io/gorm"

type NodeType string

const Primary NodeType = "primary"
const Secondary NodeType = "secondary"

type NodeStatus string

const Init NodeStatus = "init"
const Ready NodeStatus = "ready"
const Running NodeStatus = "running"
const Down NodeStatus = "down"
const Failed NodeStatus = "failed"
const Unknown NodeStatus = "unknown"

type NodeMode string

const Standalone NodeMode = "standalone"
const Cluster NodeMode = "cluster"

type Node struct {
	gorm.Model
	Type   NodeType
	Name   string `gorm:"uniqueIndex"`
	Status NodeStatus
	Mode   NodeMode
	Addr   string
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
