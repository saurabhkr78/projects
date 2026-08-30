package repository

import (
	"context"
	"errors"

	"github.com/saurabhkr78/redis-recent-views/db"
)

type FakeProductRepository struct {
	products []db.Product
}

func NewFakeProductRepository(
	products []db.Product,
) ProductRepository {
	return &FakeProductRepository{
		products: products,
	}
}

func (r *FakeProductRepository) GetProductByID(
	ctx context.Context,
	productID string,
) (*db.Product, error) {

	for i := range r.products {
		if r.products[i].ID == productID {
			return &r.products[i], nil
		}
	}

	return nil, errors.New("product not found")
}
