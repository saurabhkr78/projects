package repository

import (
	"context"
	"errors"

	"github.com/saurabhkr78/redis-recent-views/db"
)

type FakeProductRepository struct {
	products []db.Product
}

func NewFakeProductRepository(products []db.Product) ProductRepository {
	return &FakeProductRepository{
		products: products,
	}
}

func (r *FakeProductRepository) GetProductByID(
	ctx context.Context,
	productID int,
) (*db.Product, error) {

}
