package controllers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBController struct {
	dbUri     string
	Db_client *gorm.DB
}

func CreateDBController() *DBController {
	uri := os.Getenv("DB_URI")
	if uri == "" {
		log.Fatal("Please set the database uri, DB_URI")
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  uri,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	return &DBController{
		dbUri:     uri,
		Db_client: db,
	}
}
func (db *DBController) GetClient() *gorm.DB {
	return db.Db_client
}
func (db *DBController) CreateInstance(model interface{}) *gorm.DB {
	return db.Db_client.Create(model)
}
func (db *DBController) UpdateInstance(model interface{}) *gorm.DB {
	return db.Db_client.Save(model)
}
func (db *DBController) Delete(model interface{}) *gorm.DB {
	return db.Db_client.Delete(model)
}
func (db *DBController) GetAll(models interface{}) *gorm.DB {
	return db.Db_client.Find(models)
}
func (db *DBController) GetObjectByID(model interface{}, id string) *gorm.DB {
	return db.Db_client.First(model, "id = ? ", id)
}
