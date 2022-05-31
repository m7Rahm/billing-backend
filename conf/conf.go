package conf

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigInterface interface {
}
type DbConfig struct {
	Server   string
	Port     uint16
	User     string
	Password string
	Database string
}

func LoadConfig() (DbConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return DbConfig{}, err
	}
	server, isPresent := os.LookupEnv("DB_HOST")
	if !isPresent {
		return DbConfig{}, fmt.Errorf("DB_HOST not found")
	}
	port, isPresent := os.LookupEnv("DB_PORT")
	if !isPresent {
		port = "1433"
	}
	user, isPresent := os.LookupEnv("DB_USER")
	if !isPresent {
		user = "sa"
	}
	password, isPresent := os.LookupEnv("DB_PWD")
	if !isPresent {
		return DbConfig{}, fmt.Errorf("DB_PWD not found")
	}
	databaseName, isPresent := os.LookupEnv("DB_NAME")
	if !isPresent {
		databaseName = "BILLING"
	}
	portInt, err := strconv.ParseInt(port, 10, 16)
	if err != nil {
		return DbConfig{}, fmt.Errorf("port is not a number")
	}
	dbConfig := DbConfig{
		Server:   server,
		Port:     uint16(portInt),
		User:     user,
		Password: password,
		Database: databaseName,
	}
	return dbConfig, nil
}
