package repository

import (
	"context"
	"redis-url-shortner/model"
)

type URLRepository interface {
	//The repository's job is to save. It doesn't need to generate or return the short URL.
	SaveURL(ctx context.Context, url model.URL) error
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}
