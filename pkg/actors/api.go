package actors

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/api"
)

func CreateAPIActor() actor.Producer {
	return func() actor.Receiver {
		return api.CreateSystemServer()
	}
}
