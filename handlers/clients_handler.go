package handlers

import (
	"billing/backend/models"
	"billing/backend/services"
	"billing/backend/services/network"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

type ClientsHandler struct {
	clientsService services.ClientsServiceInterface
	networkService network.NetworkServiceInterface
}

func NewClientsHandler(clientsService services.ClientsServiceInterface, networkService network.NetworkServiceInterface) Handler {
	return &ClientsHandler{clientsService: clientsService, networkService: networkService}
}
func (c *ClientsHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/clients", c.GetClients)
	app.Get("/addresses", c.GetAddresses)
	app.Get("/address-types", c.GetAddressTypes)
	app.Get("/switches", c.GetSwitches)
	app.Post("/client", c.AddNewClient)
}
func (ch *ClientsHandler) GetAddresses(c *fiber.Ctx) error {
	return c.JSON(ch.clientsService.GetAddresses())
}
func (ch *ClientsHandler) AddNewClient(c *fiber.Ctx) error {
	var client models.Client
	err := c.BodyParser(&client)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Error parsing body")
	}
	clientId, err := ch.clientsService.AddNewClient(&client)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.Status(202).JSON(map[string]uint{
		"client_id": clientId,
	})
}
func (ch *ClientsHandler) GetAddressTypes(c *fiber.Ctx) error {
	return c.JSON(ch.clientsService.GetAddressTypes())
}
func (ch *ClientsHandler) GetSwitches(c *fiber.Ctx) error {
	var queryStruct struct {
		City_id       *int `query:"city" json:"city_id"`
		District_id   *int `query:"district" json:"district_id"`
		Street_id     *int `query:"street" json:"street_id"`
		Settlement_id *int `query:"settlement" json:"settlement_id"`
		Building_id   *int `query:"building" json:"building_id"`
	}
	if err := c.QueryParser(&queryStruct); err != nil {
		return c.Status(500).SendString("Error parsing query")
	}
	marshal, _ := json.Marshal(queryStruct)
	var queryMap map[string]interface{}
	if err := json.Unmarshal(marshal, &queryMap); err != nil {
		return c.Status(500).SendString("Error parsing query")
	}
	var parsedQuery map[string]interface{} = make(map[string]interface{})
	for k, v := range queryMap {
		if v != nil {
			parsedQuery[k] = v
		}
	}
	switches, err := ch.networkService.GetSwitches(parsedQuery)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(switches)
}
func (ch *ClientsHandler) GetClients(c *fiber.Ctx) error {
	return c.JSON(ch.clientsService.GetClients())
}
