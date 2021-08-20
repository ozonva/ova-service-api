package models

import (
	"fmt"
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
	got, err := NewService(1, "description", "name", "address", &tomorrowLocal)

	require.NoError(t, err, "No error should be returned for valid arguments")
	require.NotNil(t, got, "Valid Service structure should be created")
	assert.NotEqual(t, uuid.Nil, got.ID, "Non-empty ID should be generated")
	assert.True(t, tomorrowLocal.Equal(*got.WhenLocal), "Local date should be set properly")
	assert.True(t, tomorrowUTC.Equal(*got.WhenUTC), "UTC date should be set properly")
}

func TestService_WhenEmptyCalendar_ShouldCreateNewServiceWithEmptyCalendar(t *testing.T) {
	got, err := NewService(1, "description", "name", "address", nil)

	require.NoError(t, err, "No error should be returned for valid arguments")
	require.NotNil(t, got, "Valid Service structure should be created")
	assert.Nil(t, got.WhenLocal, "Local date should be nil")
	assert.Nil(t, got.WhenUTC, "UTC date should be nil")
}

func TestService_WhenEmptyUserID_ShouldReturnError(t *testing.T) {
	_, err := NewService(0, "description", "name", "address", nil)

	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "can't create service entry for non-existing user",
		err.Error(), "Incorrect error message")
}

func TestService_WhenProvidedDateInThePast_ShouldReturnError(t *testing.T) {
	got := Service{
		ID: uuid.New(),
	}

	err := got.UpdateCalendar(&yesterdayLocal)

	require.Errorf(t, err, "Error should be returned")
	assert.Equal(t, "can't update calendar to the date in the past",
		err.Error(), "Incorrect error message")
}

func TestService_ShouldBeAbleToPrintItself(t *testing.T) {
	service := Service{
		ID:             uuid.MustParse("85ae1287-9dc8-4e4a-9214-35c3debbdfba"),
		UserID:         1,
		Description:    "TO 123",
		ServiceName:    "Best TO ever",
		ServiceAddress: "In the middle of nowhere",
		WhenLocal:      &tomorrowLocal,
		WhenUTC:        &tomorrowUTC,
	}

	got := service.String()

	expected := fmt.Sprintf(`Service entry:
	ID: 			85ae1287-9dc8-4e4a-9214-35c3debbdfba
	UserID: 		1
	Description: 	TO 123
	ServiceName: 	Best TO ever
	ServiceAddress: In the middle of nowhere
	WhenLocal: 		%v
	WhenUTC: 		%v
`, service.WhenLocal, service.WhenUTC)
	assert.Equal(t, expected, got, "Service printed itself in the wrong format")
}
