package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) ListServicesV1(_ context.Context, req *pb.ListServicesV1Request) (*pb.ListServicesV1Response, error) {
	log.Info().Msg("ListServiceV1 is called...")

	if req == nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is nil")
		log.Err(invalidArgErr).Msg("Error occurred in ListServicesV1")
		return nil, invalidArgErr
	}

	services, repoErr := s.repo.ListServices(req.Limit, req.Offset)

	if repoErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred during list services: %s", repoErr.Error())
	}

	infos := make([]*pb.ServiceShortInfoV1Response, len(services))

	for i, service := range services {
		infos[i] = mapServiceToServiceShortInfoV1Response(service)
	}

	return &pb.ListServicesV1Response{
		ServiceShortInfo: infos,
	}, nil
}

func mapServiceToServiceShortInfoV1Response(service models.Service) *pb.ServiceShortInfoV1Response {
	ts := timestamppb.New(*service.WhenLocal)

	return &pb.ServiceShortInfoV1Response{
		ServiceId:   service.ID.String(),
		UserId:      service.UserID,
		ServiceName: service.ServiceName,
		When:        ts,
	}
}
