package models

type NodeMessage struct {
	FromPeerID string
	Type       NodeMessageType
}

func (m *NodeMessage) IsAliveMessage() bool {
	return m.Type == ALIVE
}

func (m *NodeMessage) IsPongMessage() bool {
	return m.Type == PONG
}

type NodeMessageType uint32

const (
	PING NodeMessageType = iota + 1
	PONG
	ELECTION
	ALIVE
	ELECTED
	OK
)
