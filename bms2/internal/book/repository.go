package book

import "context"

type Repository interface {
	Create(ctx context.Context, book *Book) error
	GetByID(ctx context.Context, id int64) (*Book, error)
	GetAll(ctx context.Context) ([]Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id int64) error
}
