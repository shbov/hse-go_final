package locationsvc

import (
	"github.com/shbov/hse-go_final/internal/location/repo"
)

type locationService struct {
	repo repo.Location
}

func ()

func New(config *service.AuthConfig, repo repo.Location) service.Location {
	return &authService{
		repo: repo,
	}
}
