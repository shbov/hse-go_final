package locationsvc

import (
	"github.com/shbov/hse-go_final/internal/location/repo"
	"github.com/shbov/hse-go_final/internal/location/service"
)

type locationService struct {
	repo repo.Location
}

func New(repo repo.Location) service.Location {
	return &locationService{
		repo: repo,
	}
}
