package repomock

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"github.com/shbov/hse-go_final/internal/driver/repo"
	"github.com/stretchr/testify/mock"
)

var _ repo.Trip = (*TripMock)(nil)

type TripMock struct {
	mock.Mock
}

func (r *TripMock) UpdateDriverIdByTripId(ctx context.Context, tripId string, userId string) error {
	args := r.Called(ctx, tripId, userId)
	return args.Error(0)
}

func (r *TripMock) ChangeTripStatus(ctx context.Context, tripId string, status trip_status.TripStatus) error {
	args := r.Called(ctx, tripId, status)
	return args.Error(0)
}

func (r *TripMock) AddTrip(ctx context.Context, trip trip.Trip) error {
	args := r.Called(ctx, trip)
	return args.Error(0)
}

func (r *TripMock) GetTripsByUserId(ctx context.Context, userId string) ([]trip.Trip, error) {
	args := r.Called(ctx, userId)
	return args.Get(0).([]trip.Trip), args.Error(1)
}

func (r *TripMock) GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*trip.Trip, error) {
	args := r.Called(ctx, userId, tripId)
	return args.Get(0).(*trip.Trip), args.Error(1)
}
