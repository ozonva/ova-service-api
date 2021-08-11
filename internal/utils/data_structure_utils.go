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

func InvertMap(originalMap map[string]int) (map[int]string, error) {
	if originalMap == nil {
		return nil, nil
	}

	result := make(map[int]string, len(originalMap))

	for key, val := range originalMap {
		if _, ok := result[val]; ok {
			return nil, fmt.Errorf("inverted key collision. Multiple keys in the original map maps to the same value \"%d\". Can't invert the map", val)
		}
		result[val] = key
	}

	return result, nil
}

func FilterSlice(slice []string, filter []string) []string {
	if slice == nil {
		return nil
	}

	var result []string

	if len(filter) == 0 {
		result = make([]string, len(slice))
		copy(result, slice)
		return result
	}

	// Convert filter slice to map to achieve effective search against filter:
	// O(n) to build filterMap and then O(1) search time complexity.
	filterMap := make(map[string]struct{}, len(filter))
	for _, filteredWord := range filter {
		filterMap[filteredWord] = struct{}{}
	}

	for _, word := range slice {
		if _, ok := filterMap[word]; ok {
			continue
		}

		result = append(result, word)
	}

	return result
}
