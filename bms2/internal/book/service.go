package book

import "context"

type BookService struct {
	repo Repository
}

func NewService(repo Repository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s *BookService) Create(ctx context.Context, req *CreateBookRequest) error {
	book := &Book{
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		PublishedAt: req.PublishedAt,
	}

	return s.repo.Create(ctx, book)
}

func (s *BookService) GetByID(ctx context.Context, id int64) (*BookResponse, error) {
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		ISBN:        book.ISBN,
		PublishedAt: book.PublishedAt,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}, nil
}

func (s *BookService) GetAll(ctx context.Context) ([]BookResponse, error) {
	books, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []BookResponse

	for _, book := range books {
		response = append(response, BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Author:      book.Author,
			ISBN:        book.ISBN,
			PublishedAt: book.PublishedAt,
			CreatedAt:   book.CreatedAt,
			UpdatedAt:   book.UpdatedAt,
		})
	}

	return response, nil
}

func (s *BookService) Update(ctx context.Context, id int64, req *UpdateBookRequest) error {
	book := &Book{
		ID:          id,
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		PublishedAt: req.PublishedAt,
	}

	return s.repo.Update(ctx, book)
}

func (s *BookService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
