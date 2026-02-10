package inmemory

import (
	"context"
	"errors"
	"testing"

	"github.com/broadcast80/ozon-task/internal/pkg/models"
)

func TestRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(r *repository)
		url       string
		alias     string
		wantErr   error
		wantURL   string
		wantAlias string
	}{
		{
			name:    "success",
			url:     "https://example.com",
			alias:   "abc123",
			wantErr: nil,
		},
		{
			name: "duplicate_alias",
			setup: func(r *repository) {
				r.Create(context.Background(), "https://example.com", "abc123")
			},
			url:     "https://new.com",
			alias:   "abc123",
			wantErr: models.ErrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(10)
			if tt.setup != nil {
				tt.setup(r)
			}

			err := r.Create(context.Background(), tt.url, tt.alias)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				gotURL, err := r.Get(context.Background(), tt.alias)
				if err != nil {
					t.Errorf("Get after Create error = %v", err)
				}
				if gotURL != tt.url {
					t.Errorf("stored URL = %q, want %q", gotURL, tt.url)
				}

				exists, _ := r.URLExists(context.Background(), tt.url)
				if !exists {
					t.Error("URLExists after Create returned false")
				}
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	r := New(10)

	r.Create(context.Background(), "https://hooli.com", "abc123")
	r.Create(context.Background(), "https://google.com", "google")

	tests := []struct {
		name    string
		alias   string
		wantURL string
		wantErr error
	}{
		{
			name:    "success",
			alias:   "abc123",
			wantURL: "https://hooli.com",
			wantErr: nil,
		},
		{
			name:    "success_other",
			alias:   "google",
			wantURL: "https://google.com",
			wantErr: nil,
		},
		{
			name:    "not_found",
			alias:   "nonexistent",
			wantErr: models.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := r.Get(context.Background(), tt.alias)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get(%q) error = %v, wantErr %v", tt.alias, err, tt.wantErr)
				return
			}
			if gotURL != tt.wantURL {
				t.Errorf("Get(%q) = %q, want %q", tt.alias, gotURL, tt.wantURL)
			}
		})
	}
}

func TestRepository_URLExists(t *testing.T) {
	r := New(10)

	r.Create(context.Background(), "https://dogville.com", "abc123")

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "exists",
			url:  "https://dogville.com",
			want: true,
		},
		{
			name: "not_exists",
			url:  "https://nonexistent.com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.URLExists(context.Background(), tt.url)
			if err != nil {
				t.Errorf("URLExists(%q) unexpected error: %v", tt.url, err)
			}
			if got != tt.want {
				t.Errorf("URLExists(%q) = %t, want %t", tt.url, got, tt.want)
			}
		})
	}
}
