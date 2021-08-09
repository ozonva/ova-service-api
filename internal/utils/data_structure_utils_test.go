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

func TestGetSliceChunks_WhenSliceIsNil_ThenShouldReturnError(t *testing.T) {
	var slice []int = nil
	chunkSize := 1

	got, err := GetSliceChunks(slice, chunkSize)

	assert.Nil(t, got, "Should not return slice when error occurs")
	assert.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "original slice doesn't exist", err.Error(), "Incorrect error message")
}

func TestGetSliceChunks_WhenChunkSizeIsZeroOrNegative_ThenShouldReturnError(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSizes := []int{0, -1}

	for _, chunkSize := range chunkSizes {
		got, err := GetSliceChunks(slice, chunkSize)

		assert.Nil(t, got, "Should not return slice when error occurs")
		assert.Errorf(t, err, "Error should be returned")
		assert.Equal(t, "chunkSize argument value must be positive", err.Error(), "Incorrect error message")
	}
}

////////////////////////////
// InvertMap tests       //
////////////////////////////
func TestInvertMap_WhenMapContainsValues_ThenShouldReturnInvertedMap(t *testing.T) {
	originalMap := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	got, err := InvertMap(originalMap)

	assert.Nil(t, err, "No error should be returned for valid original map")
	assert.Equal(t, 2, len(got), "Inverted map length should be equal original map length")
	assert.Contains(t, got, 1, "Inverted map should contains key = 1")
	assert.Equal(t, got[1], "foo", "Inverted map should keep value \"foo\" for key = 1")
	assert.Contains(t, got, 2, "Inverted map should contains key = 2")
	assert.Equal(t, got[2], "bar", "Inverted map should keep value \"bar\" for key = 2")
}

func TestInvertMap_WhenMapContainsZeroElements_ThenShouldReturnEmptyMap(t *testing.T) {
	originalMap := make(map[string]int)

	got, err := InvertMap(originalMap)

	assert.Nil(t, err, "No error should be returned for empty original map")
	assert.Empty(t, got, "Inverted map should be empty")
}

func TestInvertMap_WhenOriginalMapIsNil_ThenShouldReturnNil(t *testing.T) {
	var originalMap map[string]int = nil

	got, err := InvertMap(originalMap)

	assert.Nil(t, err, "No error should be returned when original map doesn't exist")
	assert.Nil(t, got, "Nil result should be returned")
}

func TestInvertMap_WhenMapContainsSameValueForMultipleKeys_ThenShouldReturnError(t *testing.T) {
	originalMap := map[string]int{
		"foo": 1,
		"bar": 2,
		"baz": 2,
	}

	got, err := InvertMap(originalMap)

	assert.Nil(t, got, "No result should be returned when error occurs")
	assert.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "inverted key collision. Multiple keys in the original map maps to the same value \"2\". Can't invert the map",
		err.Error(), "Incorrect error message")

}

////////////////////////////
// FilterSlice tests      //
////////////////////////////
func TestFilterSlice_WhenOriginalSliceHaveSomethingToFilter_ThenShouldReturnNewSliceWithoutFilteredValues(t *testing.T) {
	slice := []string{"bar", "baz", "foo"}
	filter := []string{"foo"}

	got := FilterSlice(slice, filter)

	assert.EqualValues(t, []string{"bar", "baz"}, got, "Result slice values are incorrect")
}

func TestFilterSlice_WhenOriginalSliceDoesNotIntersectWithFilterValues_ThenShouldReturnCopyOfTheSlice(t *testing.T) {
	slice := []string{"bar", "baz"}
	filter := []string{"foo"}

	got := FilterSlice(slice, filter)

	assert.False(t, &slice == &got, "Copy of the slice should be returned, not the original slice")
	assert.EqualValues(t, []string{"bar", "baz"}, got, "Original and result slice should contain the same values")
}

func TestFilterSlice_WhenOriginalSliceIsNil_ThenShouldReturnNil(t *testing.T) {
	var slice []string = nil
	filter := []string{"foo"}

	got := FilterSlice(slice, filter)

	assert.Nil(t, got, "Nil result should be returned")
}

func TestFilterSlice_WhenFilterIsNilOrEmpty_ThenShouldReturnCopyOfTheSlice(t *testing.T) {
	slice := []string{"bar", "baz"}
	filters := [][]string{
		nil,
		make([]string, 0),
	}

	for _, filter := range filters {
		got := FilterSlice(slice, filter)

		assert.False(t, &slice == &got, "Copy of the slice should be returned, not the original slice")
		assert.EqualValues(t, []string{"bar", "baz"}, got, "Original and result slice should contain the same values")
	}
}
