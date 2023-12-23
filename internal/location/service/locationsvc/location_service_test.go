package locationsvc

import (
	"context"
	"github.com/shbov/hse-go_final/internal/location/model"
	"github.com/shbov/hse-go_final/internal/location/repo/repomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSetLocationByDriverId(t *testing.T) {
	mockRepo := new(repomock.LocationMock)

	// Expect SetLocationByDriverId to be called with specific arguments and return no error.
	mockRepo.On("SetLocationByDriverId", mock.Anything, "driver123", float64(1.0), float64(2.0)).Return(nil)

	locationService := New(context.Background(), mockRepo)

	err := locationService.SetLocationByDriverId(context.Background(), "driver123", 1.0, 2.0)
	assert.NoError(t, err, "Expected no error")

	// Verify that the expected method was called with the expected arguments.
	mockRepo.AssertExpectations(t)
}

func TestGetDriversInLocation(t *testing.T) {
	mockRepo := new(repomock.LocationMock)

	// Expect GetDriversInLocation to be called with specific arguments and return a mock result.
	mockRepo.On("GetDriversInLocation", mock.Anything, float64(1.0), float64(2.0), float64(3.0)).
		Return([]model.Location{{Lat: 1.0, Lng: 2.0}}, nil)

	locationService := New(context.Background(), mockRepo)

	result, err := locationService.GetDriversInLocation(context.Background(), float64(1.0), float64(2.0), float64(3.0))

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected result not to be nil")
	assert.Equal(t, float64(1.0), result[0].Lat)
	assert.Equal(t, float64(2.0), result[0].Lng)

	// Verify that the expected method was called with the expected arguments.
	mockRepo.AssertExpectations(t)
}

func TestNew(t *testing.T) {
	mockRepo := new(repomock.LocationMock)

	locationService := New(context.Background(), mockRepo)

	assert.NotNil(t, locationService, "Expected service not to be nil")

	// Verify that the expected method was called with the expected arguments.
	mockRepo.AssertExpectations(t)
}
