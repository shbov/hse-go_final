package locationsvc

import (
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/location/model"
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

func (ls *locationService) GetDriversInLocation(ctx context.Context, centerLat float64, centerLng float64, radius float64) ([]model.Location, error) {
	result, err := ls.repo.GetDriversInLocation(ctx, centerLat, centerLng, radius)
	return result, err
}

func New(ctx context.Context, repo repo.Location) service.Location {
	lg := zapctx.Logger(ctx)
	s := &locationService{
		repo: repo,
	}

	lg.Info("service successfully created")
	return s
}
