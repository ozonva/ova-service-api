package api

import (
	"context"
	"github.com/ozonva/ova-service-api/internal/events"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) CreateServiceV1(_ context.Context, req *pb.CreateServiceV1Request) (*pb.CreateServiceV1Response, error) {
	log.Info().Msg("CreateServiceV1 is called...")

	if req == nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is nil")
		log.Err(invalidArgErr).Msg("Error occurred in CreateServiceV1")
		return nil, invalidArgErr
	}

	when := extractTimeFromTimestamp(req.GetWhen())
	service, err := models.NewService(req.UserId, req.Description, req.ServiceName, req.ServiceAddress, when)

	if err != nil {
		internalErr := status.Errorf(codes.Internal, "Error occurred during service creation: %s", err.Error())
		log.Err(internalErr).Msg("Error occurred in CreateServiceV1")
		return nil, internalErr
	}

	saverErr := s.saver.Save(*service)
	if saverErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred while saver trying to save the service: %s", saverErr.Error())
	}

	event := events.NewServiceCreateEvent(service.ID)
	kafkaErr := s.producer.SendMessage(event.String())
	if kafkaErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred while trying to produce Create event to Kafka: %s", kafkaErr.Error())
	}

	s.metrics.IncrementCreateCounter()

	return &pb.CreateServiceV1Response{ServiceId: service.ID.String()}, nil
}

func extractTimeFromTimestamp(ts *timestamp.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	localTime := time.Unix(ts.Seconds, int64(ts.Nanos))
	return &localTime
}
