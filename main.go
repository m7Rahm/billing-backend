package main

import (
	"billing/backend/database"
	"billing/backend/database/repository"
	"billing/backend/handlers"
	"billing/backend/models"
	"billing/backend/server"
	"billing/backend/services"
	"billing/backend/services/network"

	"github.com/gofiber/fiber/v2"
	// _ "github.com/microsoft/go-mssqldb"
)

func main() {
	db := database.GetConnection()
	db.AutoMigrate(&models.Switch{}, &models.Client{}, &models.Address{}, &models.AddressType{})
	clientRepo := repository.NewClientRepo(db)
	netRepo := repository.NewNetRepo(db)
	clientsService := services.NewClientsService(clientRepo)
	networkService := network.NewNetworkService(netRepo)
	handlers := []handlers.Handler{
		handlers.NewClientsHandler(clientsService, networkService),
		handlers.NewNetworkHandler(networkService),
	}
	// show lldp info remote-device
	app := fiber.New()
	server.RegisterRoutes(app, handlers)
	app.Listen(":4000")
}
