package api

import (
	"context"
	"github.com/ozonva/ova-service-api/internal/events"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) UpdateServiceV1(_ context.Context, req *pb.UpdateServiceV1Request) (*empty.Empty, error) {
	log.Info().Msg("UpdateServiceV1 is called...")

	if req == nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is nil")
		log.Err(invalidArgErr).Msg("Error occurred in UpdateServiceV1")
		return nil, invalidArgErr
	}

	serviceID, err := uuid.Parse(req.ServiceId)

	if err != nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is not valid UUID")
		log.Err(invalidArgErr).Msg("Error occurred in RemoveServicesV1")
		return nil, invalidArgErr
	}

	when := extractTimeFromTimestamp(req.GetWhen())
	updatedService, err := models.NewService(req.UserId, req.Description, req.ServiceName, req.ServiceAddress, when)

	if err != nil {
		internalErr := status.Errorf(codes.Internal, "Error occurred during domain service creation: %s", err.Error())
		log.Err(internalErr).Msg("Error occurred in UpdateServiceV1")
		return nil, internalErr
	}

	updatedService.ID = serviceID
	repoErr := s.repo.UpdateService(updatedService)
	if repoErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred during saving to repo: %s", repoErr.Error())
	}

	event := events.NewServiceUpdateEvent(serviceID)
	kafkaErr := s.producer.SendMessage(event.String())
	if kafkaErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred while trying to produce Update event to Kafka: %s", kafkaErr.Error())
	}

	return &empty.Empty{}, nil
}
