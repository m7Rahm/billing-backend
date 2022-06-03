package handlers

import (
	"billing/backend/models"
	"billing/backend/services/network"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type NetworkHandler struct {
	networkService network.NetworkServiceInterface
}

func NewNetworkHandler(networkService network.NetworkServiceInterface) Handler {
	return &NetworkHandler{networkService: networkService}
}
func (nh *NetworkHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/net/switch", nh.GetSwitches)
	app.Post("/net/switch", nh.AddNewSwitch)
	app.Get("/net/vendor", nh.GetVendors)
	app.Put("/net/switch/:switch/:port", nh.SetPortStatus)
	app.Get("/net/switch-info/:id", nh.GetSwitchDetails)
	app.Get("/net/switch/:id", nh.GetSwitchInfo)
}
func (nh *NetworkHandler) GetVendors(c *fiber.Ctx) error {
	vendors, err := nh.networkService.GetVendors()
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(vendors)
}
func (nh *NetworkHandler) GetSwitchDetails(c *fiber.Ctx) error {
	switchId := c.Params("id")
	switches, err := nh.networkService.GetSwitchDetails(switchId)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(switches)
}
func (nh *NetworkHandler) SetPortStatus(c *fiber.Ctx) error {
	var payload struct {
		Status string `json:"status"`
	}
	portId := c.Params("port")
	switchId := c.Params("switch")
	body := c.Body()
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	err = nh.networkService.SetPortStatus(switchId, portId, payload.Status)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(nil)
}
func (nh *NetworkHandler) AddNewSwitch(c *fiber.Ctx) error {
	body := c.Body()
	var switchInfo models.Switch
	err := json.Unmarshal(body, &switchInfo)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(err)
	}
	err = nh.networkService.AddNewSwitch(&switchInfo)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(err)
	}
	return c.SendStatus(201)
}
func (nh *NetworkHandler) GetSwitchInfo(c *fiber.Ctx) error {
	switchId := c.Params("id")
	ctx := c.Context()
	switchPort, err := nh.networkService.GetSwitchInfo(ctx, switchId)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(map[string]string{
			"error": "cannot get switch info",
		})
	}
	return c.JSON(switchPort)
}
func (nh *NetworkHandler) GetSwitches(c *fiber.Ctx) error {
	referer := c.Query("referer")
	if referer == "net" {
		name := c.Query("name")
		ip := c.Query("ip")
		vendor := c.Query("vendor")
		mac := c.Query("mac")
		switches, err := nh.networkService.GetSwitchList(name, vendor, mac, ip)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		return c.JSON(map[string]interface{}{
			"switches": switches,
		})
	} else {
		var query map[string]interface{} = make(map[string]interface{}, 0)
		city := c.Query("city")
		street := c.Query("street")
		district := c.Query("district")
		building := c.Query("building")
		if city != "" {
			query["city_id"] = city
		}
		if street != "" {
			query["street_id"] = street
		}
		if district != "" {
			query["district_id"] = district
		}
		if building != "" {
			query["building_id"] = building
		}
		switches, err := nh.networkService.GetSwitches(query)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		return c.JSON(switches)
	}
}
