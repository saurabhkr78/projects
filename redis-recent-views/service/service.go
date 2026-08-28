package service

import (
	"context"

	"github.com/saurabhkr78/redis-recent-views/db"
	"github.com/saurabhkr78/redis-recent-views/repository"
)

type ProductService interface {
	GetProduct(ctx context.Context, productID int) (*db.Product, error)
	GetRecentViews(ctx context.Context, userID string, limit int) ([]string, error)
	SetRecentView(ctx context.Context, userID, productID string) error
}

type productService struct {
	productRepo    repository.ProductRepository
	recentViewRepo repository.RecentViewRepository
}

func NewProductService(
	pRepo repository.ProductRepository,
	rRepo repository.RecentViewRepository,
) ProductService {
	return &productService{
		productRepo:    pRepo,
		recentViewRepo: rRepo,
	}
}

func (s *productService) GetProduct(
	ctx context.Context,
	productID int,
) (*db.Product, error) {
	// TODO:
	// call productRepo.GetProductByID()
	return nil, nil
}

func (s *productService) GetRecentViews(
	ctx context.Context,
	userID string,
	limit int,
) ([]string, error) {
	// TODO:
	// call recentViewRepo.GetRecentViews()
	return nil, nil
}

func (s *productService) SetRecentView(
	ctx context.Context,
	userID string,
	productID string,
) error {
	// TODO:
	// call recentViewRepo.AddRecentView()
	return nil
}
