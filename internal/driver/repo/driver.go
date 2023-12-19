package repo

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model"
)

type DriverRepo interface {
	SaveDriver(ctx context.Context, data model.Driver) error
	GetDriverById(ctx context.Context, Id string) (model.Driver, error)
}
