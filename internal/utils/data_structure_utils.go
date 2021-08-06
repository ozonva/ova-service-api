package utils

import (
	"fmt"
	"math"
)

func GetSliceChunks(slice []int, chunkSize int) ([][]int, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunkSize argument value must be positive")
	}

	n := len(slice)
	if n == 0 {
		return make([][]int, 0), nil
	}

	chunkCount := int(math.Ceil(float64(n) / float64(chunkSize)))
	result := make([][]int, chunkCount)
	k := 0

	for i := 0; i < n; i += chunkSize {
		size := i + chunkSize

		if size <= n {
			result[k] = slice[i:size]
		} else {
			result[k] = slice[i:]
		}

		k++
	}

	return result, nil
}
