package repository

import (
	"context"

	"github.com/saurabhkr78/redis-recent-views/db"
)

type ProductRepository interface {
	GetProductByID(
		ctx context.Context,
		productID string,
	) (*db.Product, error)
}

type RecentViewRepository interface {
	GetRecentViews(
		ctx context.Context,
		userID string,
		limit int,
	) ([]string, error)

	AddRecentView(
		ctx context.Context,
		userID string,
		productID string,
	) error
}
