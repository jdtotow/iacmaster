package models

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/protos/github.com/jdtotow/iacmaster/pkg/msg"
)

type Publisher struct {
	ActorEngine *actor.Engine
}

func NewPublisher(actor *actor.Engine) *Publisher {
	return &Publisher{
		ActorEngine: actor,
	}
}

func (p *Publisher) Publish(event *msg.Event, subscription *Subscription) {
	if subscription.ActionType == WebhookAction {
		//send to http post request
		url := subscription.Destination
		data := event.Data
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	} else if subscription.ActionType == ActorEngineAction {
		address, pid := strings.Split(subscription.Destination, ":")[0], strings.Split(subscription.Destination, ":")[1] // 192.168.0.4:iacmaster/system
		destinationPID := actor.NewPID(address, pid)
		p.ActorEngine.Send(destinationPID, event.Data)
	} else {
		log.Println("Unknown action type")
	}
}
