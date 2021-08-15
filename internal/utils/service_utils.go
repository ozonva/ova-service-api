package utils

import (
	"fmt"
	"github.com/ozonva/ova-service-api/internal/domain"
	"math"
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
