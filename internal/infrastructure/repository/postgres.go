package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/broadcast80/ozon-task/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	client *pgxpool.Pool
}

// func NewRepository(client *pgxpool.Pool)

func New(client *pgxpool.Pool) *repository {
	return &repository{client: client}
}

func (r *repository) Create(ctx context.Context, link domain.Link) error {
	q := `
		INSERT INTO link (url, alias, created_at) 
		VALUES ($1, $2, $3)
	`

	_, err := r.client.Exec(ctx, q, link.URL, link.Alias, link.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
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

	err := row.Scan(url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
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
