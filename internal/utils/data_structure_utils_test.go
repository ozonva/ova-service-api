package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

////////////////////////////
// GetSliceChunks tests   //
////////////////////////////
func TestGetSliceChunks_WhenLengthIsMultipleBySize_ThenShouldReturnEqualChunks(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	chunkSize := 2

	got, err := GetSliceChunks(slice, chunkSize)

	assert.Nil(t, err, "No error should be returned for positive chunkSize")
	assert.Equal(t, 2, len(got), "The result slice should contain two items")
	assert.EqualValues(t, []int{1, 2}, got[0], "The first chunk contains incorrect values")
	assert.EqualValues(t, []int{3, 4}, got[1], "The second chunk contains incorrect values")
}

func TestGetSliceChunks_WhenLengthIsNotMultipleBySize_ThenLastChunkShouldBeTheSmallest(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	chunkSize := 3

	got, err := GetSliceChunks(slice, chunkSize)

	assert.Nil(t, err, "No error should be returned for positive chunkSize")
	assert.Equal(t, len(got), 2, "The result slice should contain two items")
	assert.EqualValues(t, []int{1, 2, 3}, got[0], "The first chunk contains incorrect values")
	assert.EqualValues(t, []int{4, 5}, got[1], "The second chunk contains incorrect values")
}

func TestGetSliceChunks_WhenSliceIsEmpty_ThenShouldReturnEmptySlice(t *testing.T) {
	slice := make([]int, 0)
	chunkSize := 1

	got, err := GetSliceChunks(slice, chunkSize)

	assert.Nil(t, err, "No error should be returned for positive chunkSize")
	assert.Equal(t, 0, len(got), "The result slice should have zero length")
}

func TestGetSliceChunks_WhenChunkSizeIsGreaterSliceLength_ThenShouldReturnSingleChunkWithInitialSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSize := 10

	got, err := GetSliceChunks(slice, chunkSize)

	assert.Nil(t, err, "No error should be returned for positive chunkSize")
	assert.Equal(t, 1, len(got), "The result slice should have single item")
	assert.EqualValues(t, []int{1, 2, 3}, got[0], "The first chunk contains incorrect values")
}

func TestGetSliceChunks_WhenChunkSizeIsZeroOrNegative_ThenShouldReturnError(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSizes := []int{0, -1}

	for _, chunkSize := range chunkSizes {
		got, err := GetSliceChunks(slice, chunkSize)
		assert.Nil(t, got, "Should not return slice when error occurs")
		assert.Errorf(t, err, "chunkSize argument value must be positive", "Error with proper message should be returned")
	}
}
