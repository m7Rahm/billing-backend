package repository

import (
	"billing/backend/models"
	// "context"

	"gorm.io/gorm"
)

type ClientsRepoInterface interface {
	GetClient(val int) (*models.Client, error)
	GetClients() []models.Client
	GetAddressTypes() []models.AddressType
	GetAddresses() []models.Address
	AddNewClient(client *models.Client) (uint, error)
}

type ClientRepo struct {
	db *gorm.DB
}

func NewClientRepo(db *gorm.DB) ClientsRepoInterface {
	return &ClientRepo{db: db}
}
func (c *ClientRepo) AddNewClient(client *models.Client) (uint, error) {
	tx := c.db.Create(client)
	return client.ID, tx.Error
}
func (c *ClientRepo) GetAddresses() []models.Address {
	var addresses []models.Address
	c.db.Find(&addresses)
	return addresses
}

func (c *ClientRepo) GetAddressTypes() []models.AddressType {
	var addressTypes []models.AddressType
	c.db.Find(&addressTypes, []int{1, 4, 5})
	return addressTypes
}
func (c *ClientRepo) GetClient(val int) (*models.Client, error) {
	return &models.Client{}, nil
}

func (c *ClientRepo) GetClients() []models.Client {
	var clients []models.Client
	c.db.Joins("JOIN switches ON switches.id = clients.switch_id").Find(&clients)
	return clients
}
