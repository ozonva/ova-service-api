package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	yesterdayLocal = time.Now().AddDate(0, 0, -1)
	tomorrowLocal  = time.Now().AddDate(0, 0, 1)
	tomorrowUTC    = tomorrowLocal.UTC()
)

func TestService_WhenValidArguments_ShouldCreateNewService(t *testing.T) {
	got, err := New(1, "description", "name", "address", &tomorrowLocal)

	assert.Nil(t, err, "No error should be returned for valid arguments")
	require.NotNil(t, got, "Valid Service structure should be created")
	assert.NotEqual(t, uuid.Nil, got.ID, "Non-empty ID should be generated")
	assert.True(t, tomorrowLocal.Equal(*got.WhenLocal), "Local date should be set properly")
	assert.True(t, tomorrowUTC.Equal(*got.WhenUTC), "UTC date should be set properly")
}

func TestService_WhenEmptyCalendar_ShouldCreateNewServiceWithEmptyCalendar(t *testing.T) {
	got, err := New(1, "description", "name", "address", nil)

	assert.Nil(t, err, "No error should be returned for valid arguments")
	require.NotNil(t, got, "Valid Service structure should be created")
	assert.Nil(t, got.WhenLocal, "Local date should be nil")
	assert.Nil(t, got.WhenUTC, "UTC date should be nil")
}

func TestService_WhenEmptyUserID_ShouldReturnError(t *testing.T) {
	got, err := New(0, "description", "name", "address", nil)

	assert.Nil(t, got, "Service should not be created on error")
	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "can't create service entry for non-existing user",
		err.Error(), "Incorrect error message")
}

func TestService_WhenProvidedDateInThePast_ShouldReturnError(t *testing.T) {
	got, err := New(1, "description", "name", "address", &yesterdayLocal)

	assert.Nil(t, got, "Service should not be created on error")
	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "can't update calendar to the date in the past",
		err.Error(), "Incorrect error message")
}
