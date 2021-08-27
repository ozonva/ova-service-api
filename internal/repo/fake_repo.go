package repo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-service-api/internal/models"
)

type FakeServiceRepo struct {
	services []models.Service
}

func NewFakeServiceRepo() *FakeServiceRepo {
	todayLocal := time.Now()
	todayUTC := todayLocal.UTC()
	tomorrowLocal := time.Now().AddDate(0, 0, 1)
	tomorrowUTC := tomorrowLocal.UTC()

	return &FakeServiceRepo{
		services: []models.Service{
			{
				ID:             uuid.New(),
				UserID:         1,
				Description:    "Fake service",
				ServiceName:    "Car service",
				ServiceAddress: "In the middle of nowhere",
				WhenLocal:      &todayLocal,
				WhenUTC:        &todayUTC,
			},
			{
				ID:             uuid.New(),
				UserID:         1,
				Description:    "Fake service",
				ServiceName:    "Panzer service",
				ServiceAddress: "Somewhere",
				WhenLocal:      &tomorrowLocal,
				WhenUTC:        &tomorrowUTC,
			},
		},
	}
}

func (repo *FakeServiceRepo) AddServices(services []models.Service) error {
	for _, service := range services {
		repo.services = append(repo.services, service)
		log.Info().Msgf("Service was added to repo. Service data: %s", service.String())
	}

	return nil
}

func (repo *FakeServiceRepo) ListServices(limit uint64, offset uint64) ([]models.Service, error) {
	log.Info().Msg("List services from the repo")

	length := uint64(len(repo.services))

	if limit == 0 || offset >= length {
		return make([]models.Service, 0), nil
	}

	end := offset + limit
	if end > length {
		end = length
	}

	return repo.services[offset:end], nil
}

func (repo *FakeServiceRepo) DescribeService(serviceID uuid.UUID) (*models.Service, error) {
	for _, service := range repo.services {
		if service.ID == serviceID {
			log.Info().Msgf("Service with ID: %s is found in the repo", serviceID.String())
			return &service, nil
		}
	}

	err := fmt.Errorf("service with ID: %s was not found in the repo", serviceID.String())
	log.Err(err).Msg("Error occurred during describe service")
	return nil, err
}

func (repo *FakeServiceRepo) RemoveService(serviceID uuid.UUID) error {
	removeIdx := -1
	for i, service := range repo.services {
		if service.ID == serviceID {
			removeIdx = i
			break
		}
	}

	if removeIdx == -1 {
		err := fmt.Errorf("service with ID: %s was not found in the repo", serviceID.String())
		log.Err(err).Msg("Error occurred during describe service")
		return err
	}

	repo.services = append(repo.services[:removeIdx], repo.services[removeIdx+1:]...)
	log.Info().Msgf("Service with ID: %s was deleted from the repo", serviceID.String())
	return nil
}
