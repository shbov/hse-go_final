package repo

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model"
	"log/slog"
)

var _ DriverRepo = (*DriverSvcImpl)(nil)

type DriverSvcImpl struct {
	log        *slog.Logger
	driverRepo DriverRepo
}

func NewStatsSvcImpl(log *slog.Logger, statsRepo DriverRepo) *DriverSvcImpl {
	return &DriverSvcImpl{log: log, driverRepo: statsRepo}
}

func (s DriverSvcImpl) GetStatisticsByShortURL(ctx context.Context, id string) (model.Driver, error) {
	return s.driverRepo.GetDriverById(ctx, id)
}
