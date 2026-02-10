package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/broadcast80/ozon-task/internal/pkg/models"
	"github.com/broadcast80/ozon-task/internal/pkg/utils"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_" // использовать хеш функцию и обрубить на 5 элементов, MD5

type RepositoryInterface interface {
	Create(ctx context.Context, url string, alias string) error
	Get(ctx context.Context, alias string) (string, error)
}

type service struct {
	repository RepositoryInterface
	logger     *slog.Logger
}

func New(repository RepositoryInterface, logger *slog.Logger) *service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) GetAlias(ctx context.Context, url string) (string, error) {

	// found, err := s.repository.ExistsURL(ctx, url)
	// if err != nil {
	// 	return "", err
	// }

	// if found {
	// 	return "", errors.New("url already exists")
	// }

	alias := utils.Encode(10, charset)

	for range 10 {
		err := s.repository.Create(ctx, url, alias)
		if errors.Is(err, models.ErrDuplicate) {
			alias = utils.Encode(10, charset)
			continue
		} else if err != nil {
			return "", err
		}

		break
	}

	return alias, nil
}

func (s *service) GetURL(ctx context.Context, alias string) (string, error) {

	url, err := s.repository.Get(ctx, alias)
	if err != nil {
		s.logger.Error(err.Error())
		return "", err
	}

	return url, nil
}
