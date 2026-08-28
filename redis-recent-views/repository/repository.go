package repository

import "context"

type ProductRepository interface {
	GetProductById(ctx context.Contextm, ProductID string)
}
type RecenetViewRepository interface {
	GetRecentViewed(ctx context.Context, userID string, limit int)
	SetRecentViewed(ctx context.Context, userID, productID string)
}
