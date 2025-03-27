package actors

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/models"
)

func CreateNodeActor() actor.Producer {
	return func() actor.Receiver {
		return models.NewNode()
	}
}
