package controllers

import (
	"log"
	"log/slog"
	"reflect"

	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/models"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type EventHub struct {
	Subscribers  []*models.Subscriber
	ActorEngine  *actor.Engine
	MaxQueueSize int
	Queue        []*msg.Event
	Publisher    *models.Publisher
}

func CreateEventHub() *EventHub {
	return &EventHub{
		Subscribers:  []*models.Subscriber{},
		MaxQueueSize: 1000,
		Queue:        []*msg.Event{},
	}
}

func (e *EventHub) SubscriberExist(subscriber *models.Subscriber) bool {
	for _, s := range e.Subscribers {
		if s.GetID() == subscriber.GetID() {
			return true
		}
	}
	return false
}

func (e *EventHub) AddSubscriber(subscriber *models.Subscriber) error {
	if e.SubscriberExist(subscriber) {
		return nil
	}
	e.Subscribers = append(e.Subscribers, subscriber)
	log.Println("Subscriber added")
	return nil

}

func (e *EventHub) RemoveSubscriber(subscriber *models.Subscriber) {
	for i, s := range e.Subscribers {
		if s.GetID() == subscriber.GetID() {
			e.Subscribers = append(e.Subscribers[:i], e.Subscribers[i+1:]...)
			break
		}
	}
}

func (e *EventHub) ProcessEvent(event *msg.Event) {
	for _, subscriber := range e.Subscribers {
		for _, subcription := range subscriber.Subscriptions {
			if subcription.EventType == event.Type {
				e.Publisher.Publish(event, subcription)
			}
		}
	}
}

func (e *EventHub) AddEvent(event *msg.Event) error {
	if len(e.Queue) >= e.MaxQueueSize {
		e.Queue = e.Queue[1:] //removing the first element
	}
	e.Queue = append(e.Queue, event)
	return nil
}

func (e *EventHub) Start() {
	e.Publisher = models.NewPublisher(e.ActorEngine)
	log.Println("EventHub logic started !")
}

func (e *EventHub) ConvertActionTypeString(actionType string) models.ActionType {
	switch actionType {
	case "webhookaction":
		return models.WebhookAction
	case "emailaction":
		return models.EmailAction
	case "slackaction":
		return models.SlackAction
	case "messagingaction":
		return models.MessagingAction
	case "actorengineaction":
		return models.ActorEngineAction
	default:
		return models.UnknownAction
	}
}

func (e *EventHub) ConvertEventTypeString(eventType string) msg.EventType {
	switch eventType {
	case "log":
		return msg.EventType_LOG
	case "deployment_start":
		return msg.EventType_DEPLOYMENT_START
	case "deployment_end":
		return msg.EventType_DEPLOYMENT_COMPLETED
	case "deployment_failed":
		return msg.EventType_DEPLOYMENT_FAILED
	default:
		return msg.EventType_UNKNOWN
	}

}

func (e *EventHub) handleSubscriptionRequest(request *msg.SubscriptionRequest) {
	subscriber := &models.Subscriber{
		ID: request.Id,
	}
	for _, sub := range request.Subcriptions {
		subscription := &models.Subscription{
			ActionType:   e.ConvertActionTypeString(sub.ActionType),
			EventType:    e.ConvertEventTypeString(sub.EventType),
			DeploymentID: sub.DeploymentID,
			Destination:  sub.Destination,
		}
		subscriber.AddSubscription(subscription)
	}
	e.AddSubscriber(subscriber)
}

func (e *EventHub) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		log.Println("EventHub actor started at -> ", ctx.Engine().Address())
		e.ActorEngine = ctx.Engine()
		e.Start()
	case actor.Stopped:
		log.Println("EventHub actor has stopped")
	case actor.Initialized:
		log.Println("EventHub actor initialized")
	case *actor.PID:
		log.Println("EventHub actor has god an ID")
	case *msg.Event:
		log.Println("An event us received")
	case *msg.SubscriptionRequest:
		e.handleSubscriptionRequest(m)
	default:
		slog.Warn("EventHub got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}
