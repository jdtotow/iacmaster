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
	for _, event := range e.Queue {
		if subscriber.IsSubscribedTo(&event.Type) {
			//send to the subscriber
		}
	}
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
		if subscriber.IsSubscribedTo(&event.Type) {
			//send to the subscriber
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
	log.Println("EventHub logic started !")
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
	default:
		slog.Warn("EventHub got unknown message", "msg", m, "type", reflect.TypeOf(m).String())
	}
}
