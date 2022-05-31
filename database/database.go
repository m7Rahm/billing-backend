package database

import (
	"billing/backend/conf"
	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	config, err := conf.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		config.Server, config.User, config.Password, config.Port, config.Database)
	db, err := gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	return db
}
