package repo

import (
	"context"

	"github.com/shbov/hse-go_final/internal/location/model"
)

type Location interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddLocation(ctx context.Context, login, password, email string) error
	GetLocation(ctx context.Context, login string) (*model.Location, error)
}
