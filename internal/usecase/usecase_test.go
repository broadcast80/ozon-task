package usecase

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/broadcast80/ozon-task/internal/pkg/models"
)

type repoMock struct {
	URLExistsFn func(ctx context.Context, url string) (bool, error)
	CreateFn    func(ctx context.Context, url, alias string) error
	GetFn       func(ctx context.Context, alias string) (string, error)

	urlExistsCalls int
	createCalls    int
	getCalls       int

	lastCreateURL   string
	lastCreateAlias string
	lastGetAlias    string
}

func (m *repoMock) URLExists(ctx context.Context, url string) (bool, error) {
	m.urlExistsCalls++
	return m.URLExistsFn(ctx, url)
}

func (m *repoMock) Create(ctx context.Context, url, alias string) error {
	m.createCalls++
	m.lastCreateURL = url
	m.lastCreateAlias = alias
	return m.CreateFn(ctx, url, alias)
}

func (m *repoMock) Get(ctx context.Context, alias string) (string, error) {
	m.getCalls++
	m.lastGetAlias = alias
	return m.GetFn(ctx, alias)
}

func testLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func Test_GetAlias_URLExistsError(t *testing.T) {
	var logBuf bytes.Buffer

	wantErr := models.ErrDuplicate // должны совпасть
	repo := &repoMock{
		URLExistsFn: func(ctx context.Context, url string) (bool, error) {
			return true, nil
		},
		CreateFn: func(ctx context.Context, url, alias string) error {
			t.Fatalf("Create must not be called when URLExists fails")
			return nil
		},
	}

	s := New(repo, testLogger(&logBuf))

	alias, err := s.GetAlias(context.Background(), "https://bmstu.com")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected err=%v, got %v", wantErr, err)
	}

	if alias != "" {
		t.Fatalf("expected empty alias, got %q", alias)
	}

	if repo.createCalls != 0 {
		t.Fatalf("URLExists calls: want 1, got %d", repo.urlExistsCalls)
	}
}

func Test_GetAlias_DuplicateAlias(t *testing.T) {
	var logBuf bytes.Buffer

	createAttempts := 0
	repo := &repoMock{
		URLExistsFn: func(ctx context.Context, url string) (bool, error) {
			return false, nil
		},
		CreateFn: func(ctx context.Context, url, alias string) error {
			createAttempts++
			if createAttempts <= 2 {
				return models.ErrDuplicate
			}
			return nil
		},
	}

	s := New(repo, testLogger(&logBuf))

	gotAlias, err := s.GetAlias(context.Background(), "https://sobaka.com")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if gotAlias == "" {
		t.Fatalf("expected non-empty alias")
	}

	if repo.createCalls != 3 {
		t.Fatalf("Create calls: want 3, got %d", repo.createCalls)
	}
}

func TestService_GetAlias_CreateError(t *testing.T) {
	var logBuf bytes.Buffer

	wantErr := errors.New("insert failed")
	repo := &repoMock{
		URLExistsFn: func(ctx context.Context, url string) (bool, error) {
			return false, nil
		},
		CreateFn: func(ctx context.Context, url, alias string) error {
			return wantErr
		},
	}

	s := New(repo, testLogger(&logBuf))

	alias, err := s.GetAlias(context.Background(), "https://what.com")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected err=%v, got %v", wantErr, err)
	}

	if alias != "" {
		t.Fatalf("expected empty alias, got %q", alias)
	}

	if logBuf.Len() == 0 {
		t.Fatalf("expected log output, got empty")
	}
}

func TestService_GetURL_Success(t *testing.T) {
	var logBuf bytes.Buffer

	repo := &repoMock{
		GetFn: func(ctx context.Context, alias string) (string, error) {
			if alias != "abc" {
				t.Fatalf("expected alias abc, got %q", alias)
			}
			return "https://lostmary.com", nil
		},
	}

	s := New(repo, testLogger(&logBuf))

	url, err := s.GetURL(context.Background(), "abc")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if url != "https://lostmary.com" {
		t.Fatalf("expected url %q, got %q", "https://lostmary.com", url)
	}
	if repo.getCalls != 1 {
		t.Fatalf("Get calls: want 1, got %d", repo.getCalls)
	}
}

func TestService_GetURL_Error(t *testing.T) {
	var logBuf bytes.Buffer

	wantErr := errors.New("not found")
	repo := &repoMock{
		GetFn: func(ctx context.Context, alias string) (string, error) {
			return "", wantErr
		},
		URLExistsFn: func(ctx context.Context, url string) (bool, error) {
			t.Fatalf("URLExists not expected in GetURL")
			return false, nil
		},
		CreateFn: func(ctx context.Context, url, alias string) error {
			t.Fatalf("Create not expected in GetURL")
			return nil
		},
	}

	s := New(repo, testLogger(&logBuf))

	url, err := s.GetURL(context.Background(), "abc")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected err=%v, got %v", wantErr, err)
	}
	if url != "" {
		t.Fatalf("expected empty url, got %q", url)
	}
	if logBuf.Len() == 0 {
		t.Fatalf("expected log output, got empty")
	}
}
