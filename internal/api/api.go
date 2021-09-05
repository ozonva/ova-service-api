package api

import (
	"github.com/google/uuid"
	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

type DelayedSaver interface {
	Save(service models.Service) error
}

type MultiCreateFlusher interface {
	Flush(services []models.Service) []models.Service
}

type KafkaProducer interface {
	SendMessage(message string) (int32, int64, error)
	SendMessages(messages []string) error
}

type Repo interface {
	AddServices(services []models.Service) error
	ListServices(limit, offset uint64) ([]models.Service, error)
	DescribeService(serviceID uuid.UUID) (*models.Service, error)
	RemoveService(serviceID uuid.UUID) error
	UpdateService(service *models.Service) error
}

type GrpcApiServer struct {
	pb.UnimplementedServiceAPIServer
	repo     Repo
	saver    DelayedSaver
	flusher  MultiCreateFlusher
	producer KafkaProducer
}

func NewGrpcApiServer(repo Repo, saver DelayedSaver, flusher MultiCreateFlusher, producer KafkaProducer) *GrpcApiServer {
	return &GrpcApiServer{
		repo:     repo,
		saver:    saver,
		flusher:  flusher,
		producer: producer,
	}
}
