package api

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/ozonva/ova-service-api/internal/events"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonva/ova-service-api/internal/models"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

func (s *GrpcApiServer) MultiCreateServiceV1(ctx context.Context, req *pb.MultiCreateServiceV1Request) (*pb.MultiCreateServiceV1Response, error) {
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

	tracer := opentracing.GlobalTracer()
	multiCreateParentSpan := tracer.StartSpan("MultiCreateServiceV1", opentracing.Tag{
		Key:   "Count",
		Value: len(services),
	})
	defer multiCreateParentSpan.Finish()

	ctx = opentracing.ContextWithSpan(ctx, multiCreateParentSpan)
	notSavedServices := s.flusher.Flush(ctx, services)

	if len(notSavedServices) > 0 {
		internalErr := status.Errorf(codes.Internal, "Can't save all services properly. %d services was discarded", len(notSavedServices))
		log.Err(internalErr).Msg("Error occurred in MultiCreateServiceV1")
		return nil, internalErr
	}

	// Actually this is not correct logic because Flusher possibly save some chunks of data to database,
	// but we will produce something to Kafka only if all chunks are saved correctly. I do not handle it for now
	// because actually overall logic of storing data in two places are handled incorrect. We may fix it with:
	// 1. Outbox pattern.
	// 2. Event sourcing.
	// 3. Global transaction, but it requires that Kafka consumers should process events in an idempotent manner.
	messages := mapServicesToMessages(services)
	kafkaErr := s.producer.SendMessages(messages)
	if kafkaErr != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred while trying to produce events to Kafka for MultiCreate operation: %s", kafkaErr.Error())
	}

	s.metrics.IncrementMultiCreateCounter()

	return &pb.MultiCreateServiceV1Response{ServiceId: mapServiceToServiceIDStrings(services)}, nil
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
			return nil, err
		}

		services[i] = *service
	}

	return services, nil
}

func mapServiceToServiceIDStrings(services []models.Service) []string {
	if len(services) == 0 {
		return make([]string, 0)
	}

	serviceIDs := make([]string, len(services))

	for i, service := range services {
		serviceIDs[i] = service.ID.String()
	}

	return serviceIDs
}

func mapServicesToMessages(services []models.Service) []string {
	if len(services) == 0 {
		return make([]string, 0)
	}

	messages := make([]string, len(services))

	for i, service := range services {
		event := events.NewServiceCreateEvent(service.ID)
		messages[i] = event.String()
	}

	return messages
}
