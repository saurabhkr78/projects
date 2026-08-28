package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepositoryImpl struct {
	client *redis.Client
}

func NewRedisRepository(
	client *redis.Client,
) RecentViewRepository {
	return &RedisRepositoryImpl{
		client: client,
	}
}

func (r *RedisRepositoryImpl) AddRecentView(
	ctx context.Context,
	userID string,
	productID string,
) error {
	// LREM
	// LPUSH
	// LTRIM

	return nil
}

func (r *RedisRepositoryImpl) GetRecentViews(
	ctx context.Context,
	userID string,
	limit int,
) ([]string, error) {
	// LRANGE

	return nil, nil
}
