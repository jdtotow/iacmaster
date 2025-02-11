package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBController struct {
	dbUri     string
	db_client *gorm.DB
}

func CreateDBController(uri string) *DBController {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  uri,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	return &DBController{
		dbUri:     uri,
		db_client: db,
	}
}

func (db *DBController) Handle(context *gin.Context, concernedObject string) {
	fmt.Println("DBController reached, object called : ", concernedObject)
}
