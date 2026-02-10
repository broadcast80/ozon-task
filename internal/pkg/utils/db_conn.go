package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/broadcast80/ozon-task/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewClient(ctx context.Context, maxAttempts int, sc config.PostgresConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database) // оч плохо. стрингбилдером
	err = retry.Do(
		func() error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			var err error
			pool, err = pgxpool.New(ctx, dsn)
			return err
		},
		retry.Attempts(uint(maxAttempts)),
		retry.Delay(5*time.Second),
		retry.Context(ctx),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", maxAttempts, err)
	}

	return pool, nil
}
