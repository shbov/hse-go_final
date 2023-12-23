package tripsvc

import (
	"context"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"github.com/shbov/hse-go_final/internal/driver/repo"
	"github.com/shbov/hse-go_final/internal/driver/service"
)

var _ service.Trip = (*tripService)(nil)

type tripService struct {
	tripRepo repo.Trip
}

func (ts *tripService) ChangeTripStatus(ctx context.Context, tripId string, status trip_status.TripStatus) error {
	return ts.tripRepo.ChangeTripStatus(ctx, tripId, status)
}

func (ts *tripService) AddTrip(ctx context.Context, trip trip.Trip) error {
	return ts.tripRepo.AddTrip(ctx, trip)
}

func (ts *tripService) GetTripsByUserId(ctx context.Context, userId string) ([]trip.Trip, error) {
	result, err := ts.tripRepo.GetTripsByUserId(ctx, userId)
	return result, err
}
func (ts *tripService) GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*trip.Trip, error) {
	t, err := ts.tripRepo.GetTripByUserIdTripId(ctx, userId, tripId)
	return t, err
}

func New(ctx context.Context, tripRepo repo.Trip) service.Trip {
	lg := zapctx.Logger(ctx)
	s := &tripService{
		tripRepo: tripRepo,
	}

	lg.Info("trip service successfully created")
	return s
}
