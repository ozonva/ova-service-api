package utils

import (
	"fmt"
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

////////////////////////////
// ServicesToMap tests    //
////////////////////////////

func TestServicesToMap_WhenValidServicesSlice_ThenShouldReturnMap(t *testing.T) {
	got, err := ServicesToMap(services[:2])

	assert.Nil(t, err, "No error should be returned for valid service slice")
	require.Equal(t, 2, len(got), "Map length should be equal services length")
	serviceZeroID := services[0].ID.String()
	assert.Contains(t, got, serviceZeroID, "Key is missed in the map")
	assert.Equal(t, got[serviceZeroID], services[0], "Value is incorrect")
	serviceOneID := services[1].ID.String()
	assert.Contains(t, got, serviceOneID, "Key is missed in the map")
	assert.Equal(t, got[serviceOneID], services[1], "Value is incorrect")
}

func TestServicesToMap_WhenServicesIsEmptyOrNil_ThenShouldReturnEmptyMap(t *testing.T) {
	emptyServices := [][]domain.Service{
		make([]domain.Service, 0),
		nil,
	}

	for _, emptyService := range emptyServices {
		got, err := ServicesToMap(emptyService)

		assert.Nil(t, err, "No error should be returned for empty service slice")
		assert.Empty(t, got, "Map should be empty")
	}
}

func TestServicesToMap_WhenServicesContainDuplicates_ThenShouldReturnError(t *testing.T) {
	duplicateServiceID := uuid.New()
	servicesWithDuplicates := []domain.Service{
		{ID: duplicateServiceID},
		{ID: duplicateServiceID},
	}

	got, err := ServicesToMap(servicesWithDuplicates)

	assert.Nil(t, got, "No result should be returned when error occurs")
	require.Errorf(t, err, "Error should be returned")
	expected := fmt.Sprintf("key collision. Service with ID \"%s\" already present in the map", duplicateServiceID.String())
	assert.Equal(t, expected, err.Error(), "Incorrect error message")
}
