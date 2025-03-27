package models

type NodeMessage struct {
	FromPeerID string
	Type       NodeMessageType
}

func (m *NodeMessage) IsAliveMessage() bool {
	return m.Type == ALIVE
}

func (m *NodeMessage) LeaderHeartbeatMessage() bool {
	return m.Type == LEADERHEARTBEAT
}

func (m *NodeMessage) IsPongMessage() bool {
	return m.Type == PONG
}

type NodeMessageType uint32

const (
	ELECTION NodeMessageType = iota + 1
	ALIVE
	ELECTED
	OK
	PING
	PONG
	LEADERHEARTBEAT
)
