package link

import (
	"context"
	"fmt"

	modellink "github.com/broadcast80/ozon-task/domain/model/link"
)

type Shortener struct {
	linkDataProvider modellink.DataProvider
}

func NewShortener(dataProvider modellink.DataProvider) *Shortener {
	return &Shortener{
		linkDataProvider: dataProvider,
	}
}

func (s *Shortener) CutLink(ctx context.Context, url string) (*modellink.Link, error) {
	alias, err := s.linkDataProvider.GetAlias(ctx, url)
	if err != nil {
		return nil, fmt.Errorf(
			"s.linkDataProvider.GetAlias: %w", err,
		)
	}

	link := &modellink.Link{
		Alias: alias,
		URL:   url,
	}

	return link, nil
}

func (s *Shortener) GetFullLink(ctx context.Context, alias string) (*modellink.Link, error) {
	url, err := s.linkDataProvider.GetURL(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf(
			"s.linkDataProvider.GetUrl: %w", err,
		)
	}

	link := &modellink.Link{
		Alias: alias,
		URL:   url,
	}

	return link, nil
}
