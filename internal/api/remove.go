package api

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) RemoveServiceV1(_ context.Context, req *ova_service_api.RemoveServiceV1Request) (*empty.Empty, error) {
	log.Info().Msg("RemoveServiceV1 is called...")

	if req == nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, nilRequestArgumentMsg)
		log.Err(invalidArgErr).Msg(removeServicesV1Err)
		return nil, invalidArgErr
	}

	serviceID, err := uuid.Parse(req.Uuid)

	if err != nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is not valid UUID")
		log.Err(invalidArgErr).Msg(removeServicesV1Err)
		return nil, invalidArgErr
	}

	repoErr := s.repo.RemoveService(serviceID)

	if repoErr != nil {
		return nil, status.Error(codes.NotFound, "Service was not found")
	}

	return &empty.Empty{}, nil
}
