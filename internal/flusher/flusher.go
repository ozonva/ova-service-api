package flusher

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"

	"github.com/ozonva/ova-service-api/internal/models"
	"github.com/ozonva/ova-service-api/internal/repo"
	"github.com/ozonva/ova-service-api/internal/utils"
)

type Flusher interface {
	Flush(ctx context.Context, services []models.Service) []models.Service
}

func New(chunkSize uint, serviceRepo repo.Repo) Flusher {
	return &flusher{
		chunkSize:   chunkSize,
		serviceRepo: serviceRepo,
	}
}

type flusher struct {
	chunkSize   uint
	serviceRepo repo.Repo
}

func (f *flusher) Flush(ctx context.Context, services []models.Service) []models.Service {
	chunks, err := utils.SplitToBulks(services, f.chunkSize)

	if err != nil {
		log.Printf("Error occurs in utils.SplitToBulks: %s\n", err.Error())
		return services
	}

	unsavedServices := make([]models.Service, 0)

	for i, chunk := range chunks {
		repoErr := f.serviceRepo.AddServices(chunk)

		if repoErr != nil {
			unsavedServices = append(unsavedServices, chunk...)
			log.Printf("Services chunk #%d wasn't saved: %s\n", i, repoErr.Error())
		}

		parentSpan := opentracing.SpanFromContext(ctx)
		if parentSpan != nil {
			chunkSpan := opentracing.StartSpan("BulkCreate",
				opentracing.ChildOf(parentSpan.Context()),
				opentracing.Tag{
					Key:   "Count",
					Value: len(chunk),
				})
			chunkSpan.Finish()
		}
	}

	if len(unsavedServices) > 0 {
		return unsavedServices
	}

	return nil
}
