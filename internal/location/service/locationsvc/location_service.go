package locationsvc

import (
	"github.com/shbov/hse-go_final/internal/location/repo"
	"github.com/shbov/hse-go_final/internal/location/service"
	"golang.org/x/net/context"
)

var _ service.Location = (*locationService)(nil)

type locationService struct {
	repo repo.Location
}

func (ls *locationService) SetLocationByDriverId(ctx context.Context, driverId string, lat float64, lng float64) error {
	err := ls.repo.SetLocationByDriverId(ctx, driverId, lat, lng)
	return err
}

func New(repo repo.Location) service.Location {
	return &locationService{
		repo: repo,
	}
}
