/* What business rules should be applied?

Responsibilities:

Validation
Business logic
Call Repository */

package book

import (
	"context"
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}
func (s *Service) CreateBook(ctx context.Context, book *Book) (*Book, error) {

	//validation as business logic
	if book.Year < 1900 {
		return nil, errors.New("invalid year")
	}

	return s.repository.CreateBook(ctx, book)
}

func (s *Service) GetBooks(ctx context.Context) ([]Book, error) {
	return s.repository.GetBooks(ctx)
}

func (s *Service) GetBookByID(ctx context.Context, id string) (*Book, error) {
	return s.repository.GetBookByID(ctx, id)
}

func (s *Service) UpdateBook(ctx context.Context, id string, book *Book) (*Book, error) {
	return s.repository.UpdateBook(ctx, id, book)
}

func (s *Service) DeleteBook(ctx context.Context, id string) error {
	return s.repository.DeleteBook(ctx, id)
}
