package repository

import (
	"context"
	"errors"
	urlModel "redis-url-shortner/model"
)

type DBRepository struct {
	DB []urlModel.URL
}

func NewDBRepository(db []urlModel.URL) *DBRepository {
	return &DBRepository{
		DB: db,
	}
}

func (d *DBRepository) SaveURL(ctx context.Context, url urlModel.URL) error {
	// Implementation for saving URL
	d.DB = append(d.DB, url)
	return nil
}

func (d *DBRepository) GetOriginalURL(ctx context.Context, shortID string) (string, error) {

	for _, url := range d.DB {
		if url.ShortID == shortID {
			return url.OriginalURL, nil
		}
	}

	return "", errors.New("url not found")
}
