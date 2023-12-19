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

type locationRepo struct {
	pgxPool *pgxpool.Pool
}

func (r *locationRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *locationRepo) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	return r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(context.WithValue(ctx, repo.CtxKeyTx, tx))
	})
}

func (r *locationRepo) AddLocation(ctx context.Context, driver_id string, lat float64, lng float64) error {
	var _, err = r.conn(ctx).Exec(ctx, `INSERT INTO locations (driver_id, lat, lng) VALUES ($1, $2, $3)`, driver_id, lat, lng)
	if err != nil {
		return err
	}

	return nil
}

func (r *locationRepo) GetLocation(ctx context.Context, center_lat float64, center_lng float64, radius float64) (*model.Location, error) {
	var location model.Location

	row := r.conn(ctx).QueryRow(
		ctx,
		`SELECT id, driver_id, lat, lng, created_at FROM locations WHERE (lat - $1) * (lat - $1) + (lng - $2) * (lng - $2) <= $3`,
		center_lat, center_lng, radius*radius)
	if err := row.Scan(&location.Id, &location.DriverId, &location.Lat, &location.Lng, location.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Didn't find drivers in that location\n")
		}
		return nil, err
	}

	return &location, nil
}

//func New(config *service.AuthConfig, pgxPool *pgxpool.Pool) (repo.User, error) {
//	r := &locationRepo{
//		pgxPool: pgxPool,
//	}
//
//	ctx := context.Background()
//
//	err := r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
//		for _, user := range config.Users {
//			if err := r.AddUser(ctx, user.Login, user.Pasword, user.Email); err != nil {
//				log.Fatal(err.Error())
//			}
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	return r, nil
//}
