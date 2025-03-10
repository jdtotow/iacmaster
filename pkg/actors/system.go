package actors

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/controllers"
	"github.com/jdtotow/iacmaster/pkg/models"
)

func CreateSystemActor(channel *chan models.HTTPMessage) actor.Producer {
	return func() actor.Receiver {
		return controllers.CreateSystem(channel)
	}
}
