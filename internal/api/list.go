package api

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) ListServicesV1(_ context.Context, _ *empty.Empty) (*pb.ListServicesV1Response, error) {
	log.Info().Msg("ListServiceV1 is called...")

	// We want to list all and satisfy the Repo interface
	services, repoErr := s.repo.ListServices(^uint64(0), 0)

	if repoErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred during list services: %s", repoErr.Error())
	}

	infos := make([]*pb.ServiceShortInfoV1Response, len(services))

	for i, service := range services {
		info, mapErr := mapServiceToServiceShortInfoV1Response(&service)

		if mapErr != nil {
			return nil, status.Errorf(codes.Internal, "can't convert domain entity \"service\" at index %d to response entity: %s", i, mapErr.Error())
		}

		infos[i] = info
	}

	return &pb.ListServicesV1Response{
		ServiceShortInfo: infos,
	}, nil
}

func mapServiceToServiceShortInfoV1Response(service *models.Service) (*pb.ServiceShortInfoV1Response, error) {
	if service == nil {
		return nil, fmt.Errorf("service is nil")
	}

	var ts *timestamppb.Timestamp
	if service.WhenLocal != nil {
		ts = timestamppb.New(*service.WhenLocal)
	}

	return &pb.ServiceShortInfoV1Response{
		ServiceId:   service.ID.String(),
		UserId:      service.UserID,
		ServiceName: service.ServiceName,
		When:        ts,
	}, nil
}
