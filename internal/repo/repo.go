package repo

import "github.com/ozonva/ova-service-api/internal/models"

type Repo interface {
	AddServices(entities []models.Service) error
	ListServices(limit, offset uint64) ([]models.Service, error)
	DescribeService(entityId uint64) (*models.Service, error)
}
