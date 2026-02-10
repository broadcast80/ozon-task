package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/broadcast80/ozon-task/internal/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	client *pgxpool.Pool
}

func New(client *pgxpool.Pool) *repository {
	return &repository{client: client}
}

func (r *repository) Create(ctx context.Context, url string, alias string) error {
	q := `
		INSERT INTO link (url, alias) 
		VALUES ($1, $2)
	`

	_, err := r.client.Exec(ctx, q, url, alias)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.ErrDuplicate
			}
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where,
				pgErr.Code,
				pgErr.SQLState(),
			)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) Get(ctx context.Context, alias string) (string, error) {
	q := `
		SELECT url
		FROM link
		WHERE alias = $1	
	`

	var url string

	row := r.client.QueryRow(ctx, q, alias)

	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where,
				pgErr.Code,
				pgErr.SQLState(),
			)
			return "", newErr
		}
		return "", err
	}

	return url, nil
}

func (r *repository) URLExists(ctx context.Context, url string) (bool, error) {
	var exists bool
	err := r.client.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM link WHERE url = $1 LIMIT 1)`,
		url,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
