package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepositoryImpl struct {
	client *redis.Client
}

func NewRedisRecentViewRepository(
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

	// Redis key for this user's recent views.
	key := "recent_views:" + userID

	// Remove the product if it already exists.
	// count = 0 means remove all matching occurrences.
	if err := r.client.LRem(
		ctx,
		key,
		0,
		productID,
	).Err(); err != nil {
		return fmt.Errorf(
			"failed to remove product from recent views: %w",
			err,
		)
	}

	// Add the product to the front of the list.
	if err := r.client.LPush(
		ctx,
		key,
		productID,
	).Err(); err != nil {
		return fmt.Errorf(
			"failed to add product to recent views: %w",
			err,
		)
	}

	// Keep only the 10 most recent products.
	if err := r.client.LTrim(
		ctx,
		key,
		0,
		9,
	).Err(); err != nil {
		return fmt.Errorf(
			"failed to trim recent views: %w",
			err,
		)
	}

	return nil
}

func (r *RedisRepositoryImpl) GetRecentViews(
	ctx context.Context,
	userID string,
	limit int,
) ([]string, error) {

	if limit <= 0 {
		return []string{}, nil
	}

	key := "recent_views:" + userID

	productIDs, err := r.client.LRange(
		ctx,
		key,
		0,
		int64(limit-1),
	).Result()

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get recent views: %w",
			err,
		)
	}

	return productIDs, nil
}
