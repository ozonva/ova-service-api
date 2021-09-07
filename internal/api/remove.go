package api

import (
	"context"
	"github.com/ozonva/ova-service-api/internal/events"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) RemoveServiceV1(_ context.Context, req *pb.RemoveServiceV1Request) (*empty.Empty, error) {
	log.Info().Msg("RemoveServiceV1 is called...")

	if req == nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is nil")
		log.Err(invalidArgErr).Msg("Error occurred in RemoveServicesV1")
		return nil, invalidArgErr
	}

	serviceID, err := uuid.Parse(req.ServiceId)

	if err != nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is not valid UUID")
		log.Err(invalidArgErr).Msg("Error occurred in RemoveServicesV1")
		return nil, invalidArgErr
	}

	repoErr := s.repo.RemoveService(serviceID)
	if repoErr != nil {
		return nil, status.Error(codes.NotFound, "Service was not found")
	}

	// It is possible situation when delete actually doesn't occur because entity was already deleted,
	// but we do not handle this situation for now and consider that our consumers are idempotent.
	event := events.NewServiceDeleteEvent(serviceID)
	kafkaErr := s.producer.SendMessage(event.String())
	if kafkaErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred while trying to produce Delete event to Kafka: %s", kafkaErr.Error())
	}

	s.metrics.IncrementRemoveCounter()

	return &empty.Empty{}, nil
}
