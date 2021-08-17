package utils

import (
	"fmt"
	"math"

	"github.com/ozonva/ova-service-api/internal/domain"
)

func SplitToBulks(services []domain.Service, batchSize uint) ([][]domain.Service, error) {
	if services == nil {
		return nil, fmt.Errorf("original slice doesn't exist")
	}

	if batchSize == 0 {
		return nil, fmt.Errorf("batchSize argument value must be positive")
	}

	sliceLen := uint(len(services))
	if sliceLen == 0 {
		return make([][]domain.Service, 0), nil
	}

	batchCount := int(math.Ceil(float64(sliceLen) / float64(batchSize)))
	result := make([][]domain.Service, batchCount)
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

func ServicesToMap(services []domain.Service) (map[string]domain.Service, error) {
	result := make(map[string]domain.Service, len(services))

	for _, service := range services {
		serviceID := service.ID.String()

		if _, ok := result[serviceID]; ok {
			return nil, fmt.Errorf("key collision. Service with ID \"%s\" already present in the map", serviceID)
		}

		result[serviceID] = service
	}

	return result, nil
}
