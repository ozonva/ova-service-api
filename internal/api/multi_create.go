package api

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) MultiCreateServiceV1(_ context.Context, req *pb.MultiCreateServiceV1Request) (*empty.Empty, error) {
	log.Info().Msg("MultiCreateServiceV1 is called...")

	if req == nil {
		invalidArgErr := status.Errorf(codes.InvalidArgument, "Request argument is nil")
		log.Err(invalidArgErr).Msg("Error occurred in MultiCreateServiceV1")
		return nil, invalidArgErr
	}

	services, err := mapServiceRequestToDomainServices(req.CreateService)

	if err != nil {
		internalErr := status.Errorf(codes.InvalidArgument, "Error occurred during parsing input: %s", err.Error())
		log.Err(internalErr).Msg("Error occurred in MultiCreateServiceV1")
		return nil, internalErr
	}

	notSavedServices := s.flusher.Flush(services)

	if len(notSavedServices) > 0 {
		internalErr := status.Errorf(codes.Internal, "Can't save all services properly. %d services was discarded", len(notSavedServices))
		log.Err(internalErr).Msg("Error occurred in MultiCreateServiceV1")
		return nil, internalErr
	}

	return &empty.Empty{}, nil
}

func mapServiceRequestToDomainServices(reqServices []*pb.CreateServiceV1Request) ([]models.Service, error) {
	if len(reqServices) == 0 {
		return nil, fmt.Errorf("empty service list")
	}

	services := make([]models.Service, len(reqServices))

	for i, rs := range reqServices {
		if rs == nil {
			return nil, fmt.Errorf("list contains empty values")
		}
		when := extractTimeFromTimestamp(rs.GetWhen())
		service, err := models.NewService(rs.UserId, rs.Description, rs.ServiceName, rs.ServiceAddress, when)

		if err != nil {
			return nil, fmt.Errorf("can't create service")
		}

		services[i] = *service
	}

	return services, nil
}
