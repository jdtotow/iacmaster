package models

import (
	"slices"

	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type ActionType uint

const (
	WebhookAction ActionType = iota + 1
	EmailAction
	SlackAction
	MessagingAction
)

type Subscriber struct {
	ID           string
	ActionType   ActionType
	EventTypes   []*msg.EventType
	DeploymentID string
	Destination  string
}

func NewSubscriber(id string, actionType ActionType, eventTypes []*msg.EventType, destination, deploymentID string) *Subscriber {
	return &Subscriber{
		ID:           id,
		ActionType:   actionType,
		EventTypes:   eventTypes,
		Destination:  destination,
		DeploymentID: deploymentID,
	}
}

func (s *Subscriber) GetID() string {
	return s.ID
}
func (s *Subscriber) GetActionType() ActionType {
	return s.ActionType
}
func (s *Subscriber) GetEventTypes() []*msg.EventType {
	return s.EventTypes
}
func (s *Subscriber) IsSubscribedTo(eventType *msg.EventType) bool {
	return slices.Contains(s.EventTypes, eventType)
}
