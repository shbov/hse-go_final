package tripsvc

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"github.com/shbov/hse-go_final/internal/driver/repo/repomock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTrip(t *testing.T) {
	mockRepo := new(repomock.TripMock)
	trip1 := trip.Trip{
		Id:       "trip1",
		DriverId: "driver1",
		From: trip.Coordinates{
			Lat: 1.0,
			Lng: 1.0,
		},
		To: trip.Coordinates{
			Lat: 2.0,
			Lng: 2.0,
		},
		Price: trip.Price{
			Amount:   100,
			Currency: "RUB",
		},
		Status: trip_status.DRIVERFOUND,
	}

	mockRepo.On("AddTrip", context.Background(), trip1).Return(nil)
	tripService := New(context.Background(), mockRepo)

	err := tripService.AddTrip(context.Background(), trip1)
	assert.NoError(t, err, "Expected no error")

	mockRepo.AssertExpectations(t)
}

func TestUpdateDriverIdByTripId(t *testing.T) {
	mockRepo := new(repomock.TripMock)
	tripID := "trip1"
	updatedDriverId := "driver2"

	mockRepo.On("UpdateDriverIdByTripId", context.Background(), tripID, updatedDriverId).Return(nil)
	tripService := New(context.Background(), mockRepo)

	err := tripService.UpdateDriverIdByTripId(context.Background(), tripID, updatedDriverId)

	assert.NoError(t, err, "Expected no error")

	mockRepo.AssertExpectations(t)
}

func TestGetTripsByUserId(t *testing.T) {
	// Setup
	mockRepo := new(repomock.TripMock)

	tripService := New(context.Background(), mockRepo)
	userID := "driver1"

	// Mock the repository method
	trips := []trip.Trip{
		{
			Id:       "trip1",
			DriverId: userID,
			From: trip.Coordinates{
				Lat: 1.0,
				Lng: 1.0,
			},
			To: trip.Coordinates{
				Lat: 2.0,
				Lng: 2.0,
			},
			Price: trip.Price{
				Amount:   100,
				Currency: "RUB",
			},
			Status: trip_status.DRIVERFOUND,
		},
		{
			Id:       "trip2",
			DriverId: userID,
			From: trip.Coordinates{
				Lat: 3.0,
				Lng: 3.0,
			},
			To: trip.Coordinates{
				Lat: 4.0,
				Lng: 4.0,
			},
			Price: trip.Price{
				Amount:   200,
				Currency: "RUB",
			},
			Status: trip_status.DRIVERFOUND,
		},
	}

	mockRepo.On("GetTripsByUserId", context.Background(), userID).Return(trips, nil)

	// Test
	result, err := tripService.GetTripsByUserId(context.Background(), userID)

	// Assertions
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected result not to be nil")
	assert.Len(t, result, 2, "Expected two trips in result")
	mockRepo.AssertExpectations(t)
}

func TestGetTripByUserIdTripId(t *testing.T) {
	// Setup
	mockRepo := new(repomock.TripMock)

	tripService := New(context.Background(), mockRepo)
	userID := "user123"
	tripID := "trip1"

	// Mock the repository method
	trip := trip.Trip{
		Id:       tripID,
		DriverId: userID,
		From: trip.Coordinates{
			Lat: 1.0,
			Lng: 1.0,
		},
		To: trip.Coordinates{
			Lat: 2.0,
			Lng: 2.0,
		},
		Price: trip.Price{
			Amount:   100,
			Currency: "RUB",
		},
		Status: trip_status.DRIVERFOUND,
	}

	mockRepo.On("GetTripByUserIdTripId", context.Background(), userID, tripID).Return(&trip, nil)

	// Test
	result, err := tripService.GetTripByUserIdTripId(context.Background(), userID, tripID)

	// Assertions
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected result not to be nil")
	assert.Equal(t, trip, *result, "Expected the same trip in result")
	mockRepo.AssertExpectations(t)
}
