package application

import (
	"context"

	"github.com/broadcast80/ozon-task/internal/domain"
	"github.com/broadcast80/ozon-task/internal/pkg/utils"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type RepositoryInterface interface {
	Create(ctx context.Context, link domain.Link) error
	Get(ctx context.Context, alias string) (string, error)
}

type service struct {
	repository RepositoryInterface
}

func New(repository RepositoryInterface) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(ctx context.Context, link domain.Link) (string, error) {

	// тут возможно стоит вынести в конфиг сайз и чарсет
	alias := utils.Encode(10, charset)

	link.Alias = alias

	err := s.repository.Create(ctx, link)
	if err != nil {
		// обработать
	}

	return alias, nil
}

func (s *service) Get(ctx context.Context, alias string) (string, error) {

	url, err := s.repository.Get(ctx, alias)
	if err != nil {
		// обработать
	}

	return url, nil
}
