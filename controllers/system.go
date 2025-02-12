package controllers

import (
	"errors"
	"log"

	"github.com/jdtotow/iacmaster/controllers"
	"github.com/jdtotow/iacmaster/models"
)

type System struct {
	node         *models.Node
	dbController *controllers.DBController
}

func CreateSystem(nodeType, nodeName string, dbController *controllers.DBController) *System {
	n := &models.Node{
		Type:   models.NodeType(nodeType),
		Name:   nodeName,
		Status: models.NodeStatus("init"),
	}
	return &System{
		node:         n,
		dbController: dbController,
	}
}
func (s *System) CheckMandatoryTableAndData() error {
	err := errors.New("Missing database initial data")
	return err
}

func (s *System) Start() {
	if s.node.Type == models.Primary {
		err := s.CheckMandatoryTableAndData()
		if err != nil {
			log.Fatal("Cannot continue, missing mandatory data")
		}
	}
}
