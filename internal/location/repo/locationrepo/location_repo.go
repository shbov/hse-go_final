package locationrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shbov/hse-go_final/internal/location/model"
	"github.com/shbov/hse-go_final/internal/location/repo"
)

var _ repo.Location = (*locationRepo)(nil)

type locationRepo struct {
	pgxPool *pgxpool.Pool
}

func (r *locationRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *locationRepo) AddLocation(ctx context.Context, driverId string, lat float64, lng float64) error {
	var _, err = r.conn(ctx).Exec(ctx, `INSERT INTO locations (driver_id, lat, lng) VALUES ($1, $2, $3)`, driverId, lat, lng)
	if err != nil {
		return err
	}

	return nil
}

func (r *locationRepo) GetLocation(ctx context.Context, centerLat float64, centerLng float64, radius float64) (*model.Location, error) {
	var location model.Location

	row := r.conn(ctx).QueryRow(
		ctx,
		`SELECT id, driver_id, lat, lng, created_at FROM locations WHERE (lat - $1) * (lat - $1) + (lng - $2) * (lng - $2) <= $3`,
		centerLat, centerLng, radius*radius)
	if err := row.Scan(&location.Id, &location.DriverId, &location.Lat, &location.Lng, &location.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Didn't find drivers in that location\n")
		}
		return nil, err
	}

	return &location, nil
}

func (r *locationRepo) SetLocationByDriverId(ctx context.Context, driverId string, lat float64, lng float64) error {
	_, err := r.conn(ctx).Exec(
		ctx,
		`UPDATE locations SET lat = $1, lng = $2 WHERE driver_id = $3`,
		lat, lng, driverId)
	if err != nil {
		return err
	}

	return nil
}

func New(pgxPool *pgxpool.Pool) (repo.Location, error) {
	r := &locationRepo{
		pgxPool: pgxPool,
	}

	return r, nil
}
