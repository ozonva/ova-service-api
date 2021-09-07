package repo

import (
	"github.com/google/uuid"

	"github.com/ozonva/ova-service-api/internal/models"
)

type Repo interface {
	AddServices(services []models.Service) error
	ListServices(limit uint64, offset uint64) ([]models.Service, error)
	DescribeService(serviceID uuid.UUID) (*models.Service, error)
	RemoveService(serviceID uuid.UUID) error
	UpdateService(service *models.Service) error
}
