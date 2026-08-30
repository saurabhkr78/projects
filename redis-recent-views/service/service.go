package service

import (
	"context"
	"fmt"

	"github.com/saurabhkr78/redis-recent-views/db"
	"github.com/saurabhkr78/redis-recent-views/repository"
)

type ProductService interface {
	GetProduct(
		ctx context.Context,
		productID string,
	) (*db.Product, error)

	GetRecentViews(
		ctx context.Context,
		userID string,
		limit int,
	) ([]string, error)

	SetRecentView(
		ctx context.Context,
		userID string,
		productID string,
	) error
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
	productID string,
) (*db.Product, error) {

	product, err := s.productRepo.GetProductByID(
		ctx,
		productID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get product: %w",
			err,
		)
	}

	return product, nil
}

func (s *productService) GetRecentViews(
	ctx context.Context,
	userID string,
	limit int,
) ([]string, error) {

	if limit <= 0 {
		return []string{}, nil
	}

	productIDs, err := s.recentViewRepo.GetRecentViews(
		ctx,
		userID,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get recent views: %w",
			err,
		)
	}

	return productIDs, nil
}

func (s *productService) SetRecentView(
	ctx context.Context,
	userID string,
	productID string,
) error {

	if err := s.recentViewRepo.AddRecentView(
		ctx,
		userID,
		productID,
	); err != nil {
		return fmt.Errorf(
			"failed to set recent view: %w",
			err,
		)
	}

	return nil
}
