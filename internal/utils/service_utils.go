package utils

import (
	"fmt"
	"math"

	"github.com/ozonva/ova-service-api/internal/models"
)

func SplitToBulks(services []models.Service, batchSize uint) ([][]models.Service, error) {
	if services == nil {
		return nil, fmt.Errorf("original slice doesn't exist")
	}

	if batchSize == 0 {
		return nil, fmt.Errorf("batchSize argument value must be positive")
	}

	sliceLen := uint(len(services))
	if sliceLen == 0 {
		return make([][]models.Service, 0), nil
	}

	batchCount := int(math.Ceil(float64(sliceLen) / float64(batchSize)))
	result := make([][]models.Service, batchCount)
	k := 0

	for i := uint(0); i < sliceLen; i += batchSize {
		size := i + batchSize

		if size < sliceLen {
			result[k] = services[i:size]
		} else {
			result[k] = services[i:]
		}

		k++
	}

	return result, nil
}

func ServicesToMap(services []models.Service) (map[string]models.Service, error) {
	result := make(map[string]models.Service, len(services))

	for _, service := range services {
		serviceID := service.ID.String()

		if _, ok := result[serviceID]; ok {
			return nil, fmt.Errorf("key collision. Service with ID \"%s\" already present in the map", serviceID)
		}

		result[serviceID] = service
	}

	return result, nil
}
