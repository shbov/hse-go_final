package locationsvc

import (
	"github.com/shbov/hse-go_final/internal/location/model"
	"github.com/shbov/hse-go_final/internal/location/repo"
	"github.com/shbov/hse-go_final/internal/location/service"
	"golang.org/x/net/context"
	"log"
)

var _ service.Location = (*locationService)(nil)

type locationService struct {
	repo repo.Location
}

func (ls *locationService) SetLocationByDriverId(ctx context.Context, driverId string, lat float32, lng float32) error {
	err := ls.repo.SetLocationByDriverId(ctx, driverId, lat, lng)
	return err
}

func (ls *locationService) GetDriversInLocation(ctx context.Context, centerLat float32, centerLng float32, radius float32) ([]model.Location, error) {
	result, err := ls.repo.GetDriversInLocation(ctx, centerLat, centerLng, radius)
	return result, err
}

func New(repo repo.Location) service.Location {
	s := &locationService{
		repo: repo,
	}

	log.Println("service successfully created")
	return s
}
