package utils

import (
	"fmt"
	"math"
)

func GetSliceChunks(slice []int, chunkSize int) ([][]int, error) {
	if slice == nil {
		return nil, fmt.Errorf("original slice doesn't exist")
	}

	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunkSize argument value must be positive")
	}

	sliceLen := len(slice)
	if sliceLen == 0 {
		return make([][]int, 0), nil
	}

	chunkCount := int(math.Ceil(float64(sliceLen) / float64(chunkSize)))
	result := make([][]int, chunkCount)
	k := 0

	for i := 0; i < sliceLen; i += chunkSize {
		size := i + chunkSize

		if size < sliceLen {
			result[k] = slice[i:size]
		} else {
			result[k] = slice[i:]
		}

		k++
	}

	return result, nil
}

func InvertMap(originalMap map[string]int) map[int]string {
	if originalMap == nil {
		return nil
	}

	result := make(map[int]string, len(originalMap))

	for key, val := range originalMap {
		result[val] = key
	}

	return result
}
