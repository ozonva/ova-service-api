package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozonva/ova-service-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var services = []domain.Service{
	{ID: uuid.New()},
	{ID: uuid.New()},
	{ID: uuid.New()},
	{ID: uuid.New()},
	{ID: uuid.New()},
}

////////////////////////////
// SplitToBulk tests   	  //
////////////////////////////
func TestSplitToBulks_WhenLenIsGreaterThanBatchSize_ThenShouldReturnSeveralBatches(t *testing.T) {
	batchSize := uint(3)

	got, err := SplitToBulks(services, batchSize)

	assert.Nil(t, err, "No error should be returned for positive batchSize")
	require.Equal(t, len(got), 2, "The result slice should contain two items")
	assert.EqualValues(t, services[:3], got[0], "The first batch contains incorrect values")
	assert.EqualValues(t, services[3:], got[1], "The second batch contains incorrect values")
}

func TestSplitToBulks_WhenSliceIsEmpty_ThenShouldReturnEmptySlice(t *testing.T) {
	slice := make([]domain.Service, 0)
	batchSize := uint(1)

	got, err := SplitToBulks(slice, batchSize)

	assert.Nil(t, err, "No error should be returned for positive batchSize")
	assert.Equal(t, 0, len(got), "The result slice should have zero length")
}

func TestSplitToBulks_WhenSliceIsNil_ThenShouldReturnError(t *testing.T) {
	var slice []domain.Service = nil
	batchSize := uint(1)

	got, err := SplitToBulks(slice, batchSize)

	assert.Nil(t, got, "Should not return slice when error occurs")
	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "original slice doesn't exist", err.Error(), "Incorrect error message")
}

func TestSplitToBulks_WhenChunkSizeIsZero_ThenShouldReturnError(t *testing.T) {
	batchSize := uint(0)

	got, err := SplitToBulks(services, batchSize)

	assert.Nil(t, got, "Should not return slice when error occurs")
	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "batchSize argument value must be positive", err.Error(), "Incorrect error message")
}
