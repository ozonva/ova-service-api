package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ozonva/ova-service-api/internal/models"
	"github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) DescribeServiceV1(_ context.Context, req *ova_service_api.DescribeServiceV1Request) (*ova_service_api.DescribeServiceV1Response, error) {
	log.Info().Msg("DescribeServiceV1 is called...")

	if req == nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, nilRequestArgumentMsg)
		log.Err(invalidArgErr).Msg(describeServiceV1Err)
		return nil, invalidArgErr
	}

	serviceID, err := uuid.Parse(req.Uuid)

	if err != nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is not valid UUID")
		log.Err(invalidArgErr).Msg(describeServiceV1Err)
		return nil, invalidArgErr
	}

	service, repoErr := s.repo.DescribeService(serviceID)

	if repoErr != nil {
		return nil, status.Error(codes.NotFound, "Service was not found")
	}

	return mapServiceToDescribeV1Response(service), nil
}

func mapServiceToDescribeV1Response(service *models.Service) *ova_service_api.DescribeServiceV1Response {
	ts := timestamppb.New(*service.WhenLocal)
	tsUTC := timestamppb.New(*service.WhenUTC)

	return &ova_service_api.DescribeServiceV1Response{
		Uuid:           service.ID.String(),
		UserId:         service.UserID,
		Description:    service.Description,
		ServiceName:    service.ServiceName,
		ServiceAddress: service.ServiceAddress,
		When:           ts,
		WhenUtc:        tsUTC,
	}
}
