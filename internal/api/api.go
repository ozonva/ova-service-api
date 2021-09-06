package api

import (
	"context"
	"github.com/google/uuid"
	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

type Metrics interface {
	IncrementCreateCounter()
	IncrementMultiCreateCounter()
	IncrementUpdateCounter()
	IncrementRemoveCounter()
}

type DelayedSaver interface {
	Save(service models.Service) error
}

type MultiCreateFlusher interface {
	Flush(ctx context.Context, services []models.Service) []models.Service
}

type KafkaProducer interface {
	SendMessage(message string) error
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
	metrics  Metrics
}

func NewGrpcApiServer(repo Repo, saver DelayedSaver, flusher MultiCreateFlusher, producer KafkaProducer, metrics Metrics) *GrpcApiServer {
	return &GrpcApiServer{
		repo:     repo,
		saver:    saver,
		flusher:  flusher,
		producer: producer,
		metrics:  metrics,
	}
}
