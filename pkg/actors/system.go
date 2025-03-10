package actors

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/controllers"
)

func CreateSystemActor() actor.Producer {
	return func() actor.Receiver {
		return controllers.CreateSystem()
	}
}
