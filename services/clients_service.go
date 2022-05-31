package services

import (
	"billing/backend/database/repository"
	"billing/backend/models"
)

type ClientsServiceInterface interface {
	GetClients() []models.Client
	GetAddresses() []models.Address
	GetAddressTypes() []models.AddressType
}
type ClientsService struct {
	clientRepo repository.ClientsRepoInterface
}

func NewClientsService(clientRepo repository.ClientsRepoInterface) *ClientsService {
	return &ClientsService{clientRepo: clientRepo}
}
func (s *ClientsService) GetAddresses() []models.Address {
	return s.clientRepo.GetAddresses()
}
func (s *ClientsService) GetAddressTypes() []models.AddressType {
	return s.clientRepo.GetAddressTypes()
}
func (c *ClientsService) GetClients() []models.Client {
	return c.clientRepo.GetClients()
}
