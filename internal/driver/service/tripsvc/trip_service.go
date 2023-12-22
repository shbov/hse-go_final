package tripsvc

import (
	"context"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/driver/model"
	"github.com/shbov/hse-go_final/internal/driver/repo"
	"github.com/shbov/hse-go_final/internal/driver/service"
)

var _ service.Trip = (*tripService)(nil)

type tripService struct {
	tripRepo repo.Trip
}

func (ts *tripService) AddTrip(ctx context.Context, trip model.Trip) error {
	return ts.AddTrip(ctx, trip)
}

func (ts *tripService) GetTripsByUserId(ctx context.Context, userId string) ([]model.Trip, error) {
	result, err := ts.tripRepo.GetTripsByUserId(ctx, userId)
	return result, err
}
func (ts *tripService) GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*model.Trip, error) {
	trip, err := ts.tripRepo.GetTripByUserIdTripId(ctx, userId, tripId)
	return trip, err
}

func New(ctx context.Context, tripRepo repo.Trip) service.Trip {
	lg := zapctx.Logger(ctx)
	s := &tripService{
		tripRepo: tripRepo,
	}

	lg.Info("trip service successfully created")
	return s
}
