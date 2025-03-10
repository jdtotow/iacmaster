package actors

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/jdtotow/iacmaster/pkg/controllers"
	"github.com/jdtotow/iacmaster/pkg/models"
)

func CreateRunnerActor(workingDir, name string, mandatory_commands []string, kind models.ExecutorKind, engine *actor.Engine) actor.Producer {
	return func() actor.Receiver {
		return controllers.CreateIaCRunner(workingDir, name, mandatory_commands, kind, engine)
	}
}
