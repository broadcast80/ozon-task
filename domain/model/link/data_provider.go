package link

import "context"

type DataProvider interface {
	GetAlias(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, alias string) (string, error)
}
