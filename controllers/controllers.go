package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type DBController struct {
	dbUri      string
	dbPort     int
	dbUsername string
	dbPassword string
	dbName     string
}

func CreateDBController(uri, username, password, dbname string, port int) *DBController {
	return &DBController{
		dbUri:      uri,
		dbPort:     port,
		dbUsername: username,
		dbPassword: password,
		dbName:     dbname,
	}
}

func (db *DBController) Connect() error {
	err := errors.New("could not connect to DB")
	return err
}

func (db *DBController) handle(context *gin.Context, concernedObject string) {

}
