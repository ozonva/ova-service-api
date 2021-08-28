package api

import (
	"github.com/google/uuid"

	"github.com/ozonva/ova-service-api/internal/models"
	"github.com/ozonva/ova-service-api/internal/repo"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

type Repo interface {
	AddServices(services []models.Service) error
	ListServices(limit, offset uint64) ([]models.Service, error)
	DescribeService(serviceID uuid.UUID) (*models.Service, error)
	RemoveService(serviceID uuid.UUID) error
}

type GrpcApiServer struct {
	pb.UnimplementedServiceAPIServer
	repo Repo
}

func NewGrpcApiServer() *GrpcApiServer {
	return &GrpcApiServer{
		repo: repo.NewFakeServiceRepo(),
	}
}