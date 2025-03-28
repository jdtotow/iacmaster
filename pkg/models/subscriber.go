package models

import (
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type ActionType uint

const (
	WebhookAction ActionType = iota + 1
	EmailAction
	SlackAction
	MessagingAction
	ActorEngineAction
	UnknownAction
)

type Subscription struct {
	ActionType   ActionType
	EventType    msg.EventType
	DeploymentID string
	Destination  string
}

type Subscriber struct {
	ID            string
	Subscriptions []*Subscription
}

func NewSubscriber(id string) *Subscriber {
	return &Subscriber{
		ID: id,
	}
}

func (s *Subscriber) IsSubscribedToType(eventType msg.EventType) bool {
	for _, sub := range s.Subscriptions {
		if sub.EventType == eventType {
			return true
		}
	}
	return false
}

func (s *Subscriber) IsSubscribedToDeployment(eventType msg.EventType, deploymentID string) bool {
	for _, sub := range s.Subscriptions {
		if sub.EventType == eventType && sub.DeploymentID == deploymentID {
			return true
		}
	}
	return false
}

func (s *Subscriber) AddSubscription(sub *Subscription) {
	if s.IsSubscribedToType(sub.EventType) {
		return
	}
	s.Subscriptions = append(s.Subscriptions, sub)
}

func (s *Subscriber) GetID() string {
	return s.ID
}
