package flusher

import (
	"log"

	"github.com/ozonva/ova-service-api/internal/models"
	"github.com/ozonva/ova-service-api/internal/repo"
	"github.com/ozonva/ova-service-api/internal/utils"
)

type Flusher interface {
	Flush(services []models.Service) []models.Service
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

func (f *flusher) Flush(services []models.Service) []models.Service {
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
	}

	if len(unsavedServices) > 0 {
		return unsavedServices
	}

	return nil
}
