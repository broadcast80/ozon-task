package postgresql

import (
	"context"
	"testing"
	"time"

	"github.com/broadcast80/ozon-task/internal/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:12-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS link (
            id SERIAL PRIMARY KEY,
            url TEXT NOT NULL,
            alias TEXT UNIQUE NOT NULL
        );
    `)
	require.NoError(t, err)

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}

	return pool, cleanup
}

func TestRepository_Create_Success(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()

	err := repo.Create(ctx, "https://example.com", "test")
	require.NoError(t, err)

	var url string
	err = pool.QueryRow(ctx, "SELECT url FROM link WHERE alias = $1", "test").Scan(&url)
	require.NoError(t, err)
	require.Equal(t, "https://example.com", url)
}

func TestRepository_Create_Duplicate(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()

	require.NoError(t, repo.Create(ctx, "https://example.com", "test"))

	err := repo.Create(ctx, "https://example2.com", "test")
	require.ErrorIs(t, err, models.ErrDuplicate)
}

func TestRepository_Get_Success(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()

	// Создаем запись
	require.NoError(t, repo.Create(ctx, "https://example.com", "test"))

	// Получаем
	url, err := repo.Get(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, "https://example.com", url)
}

func TestRepository_Get_NotFound(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()

	url, err := repo.Get(ctx, "nonexistent")
	require.ErrorIs(t, err, models.ErrNotFound)
	require.Empty(t, url)
}

func TestRepository_URLExists_True(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()

	require.NoError(t, repo.Create(ctx, "https://example.com", "test"))

	exists, err := repo.URLExists(ctx, "https://example.com")
	require.NoError(t, err)
	require.True(t, exists)
}

func TestRepository_URLExists_False(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()

	exists, err := repo.URLExists(ctx, "https://nonexistent.com")
	require.NoError(t, err)
	require.False(t, exists)
}
