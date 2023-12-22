package repomock

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model"
	"github.com/shbov/hse-go_final/internal/driver/repo"
	"github.com/stretchr/testify/mock"
)

var _ repo.Trip = (*TripMock)(nil)

type TripMock struct {
	mock.Mock
}

func (r *TripMock) AddTrip(ctx context.Context, trip model.Trip) error {
	args := r.Called(ctx, trip)
	return args.Error(0)
}

func (r *TripMock) GetTripsByUserId(ctx context.Context, userId string) ([]model.Trip, error) {
	args := r.Called(ctx, userId)
	return args.Get(0).([]model.Trip), args.Error(1)
}

func (r *TripMock) GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*model.Trip, error) {
	args := r.Called(ctx, userId, tripId)
	return args.Get(0).(*model.Trip), args.Error(1)
}
